package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/Dreamacro/clash/adapters/outbound"
	"github.com/Dreamacro/clash/adapters/outboundgroup"
	"github.com/Dreamacro/clash/adapters/provider"
	"github.com/Dreamacro/clash/component/auth"
	"github.com/Dreamacro/clash/component/fakeip"
	"github.com/Dreamacro/clash/component/trie"
	CC "github.com/Dreamacro/clash/config"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/dns"
	"github.com/Dreamacro/clash/log"
	R "github.com/Dreamacro/clash/rules"
	T "github.com/Dreamacro/clash/tunnel"
)

// DNS config
type DNS struct {
	Enable            bool             `json:"enable"`
	IPv6              bool             `json:"ipv6"`
	NameServer        []dns.NameServer `json:"nameserver"`
	Fallback          []dns.NameServer `json:"fallback"`
	FallbackFilter    FallbackFilter   `json:"fallback-filter"`
	Listen            string           `json:"listen"`
	EnhancedMode      dns.EnhancedMode `json:"enhanced-mode"`
	DefaultNameserver []dns.NameServer `json:"default-nameserver"`
	FakeIPRange       *fakeip.Pool
	Hosts             *trie.DomainTrie
}

// FallbackFilter config
type FallbackFilter struct {
	GeoIP  bool         `json:"geoip"`
	IPCIDR []*net.IPNet `json:"ipcidr"`
	Domain []string     `json:"domain"`
}

type RawDNS struct {
	Enable            bool              `json:"enable"`
	IPv6              bool              `json:"ipv6"`
	UseHosts          bool              `json:"use-hosts"`
	NameServer        []string          `json:"nameserver"`
	Fallback          []string          `json:"fallback"`
	FallbackFilter    RawFallbackFilter `json:"fallback-filter"`
	Listen            string            `json:"listen"`
	EnhancedMode      dns.EnhancedMode  `json:"enhanced-mode"`
	FakeIPRange       string            `json:"fake-ip-range"`
	FakeIPFilter      []string          `json:"fake-ip-filter"`
	DefaultNameserver []string          `json:"default-nameserver"`
}

type RawFallbackFilter struct {
	GeoIP  bool     `json:"geoip"`
	IPCIDR []string `json:"ipcidr"`
	Domain []string `json:"domain"`
}

type RawConfig struct {
	Port               int          `json:"port"`
	SocksPort          int          `json:"socks-port"`
	RedirPort          int          `json:"redir-port"`
	MixedPort          int          `json:"mixed-port"`
	Authentication     []string     `json:"authentication"`
	AllowLan           bool         `json:"allow-lan"`
	BindAddress        string       `json:"bind-address"`
	Mode               T.TunnelMode `json:"mode"`
	LogLevel           log.LogLevel `json:"log-level"`
	IPv6               bool         `json:"ipv6"`
	ExternalController string       `json:"external-controller"`
	ExternalUI         string       `json:"external-ui"`
	Secret             string       `json:"secret"`
	Interface          string       `json:"interface-name"`

	ProxyProvider map[string]map[string]interface{} `json:"proxy-providers"`
	Hosts         map[string]string                 `json:"hosts"`
	DNS           RawDNS                            `json:"dns"`
	Experimental  CC.Experimental                   `json:"experimental"`
	Proxy         []map[string]interface{}          `json:"proxies"`
	ProxyGroup    []map[string]interface{}          `json:"proxy-groups"`
	Rule          []string                          `json:"rules"`
}

// Parse config
func Parse(profile string, cfg string) (*CC.Config, error) {
	rawCfg, err := UnmarshalRawConfig(profile, cfg)
	if err != nil {
		return nil, err
	}

	return ParseRawConfig(rawCfg)
}

func UnmarshalRawConfig(profile string, cfg string) (*RawConfig, error) {
	// config with some default value
	rawCfg := &RawConfig{
		AllowLan:       false,
		BindAddress:    "*",
		Mode:           T.Rule,
		Authentication: []string{},
		LogLevel:       log.INFO,
		Hosts:          map[string]string{},
		Rule:           []string{},
		Proxy:          []map[string]interface{}{},
		ProxyGroup:     []map[string]interface{}{},
		DNS: RawDNS{
			Enable:      false,
			UseHosts:    true,
			FakeIPRange: "198.18.0.1/16",
			FallbackFilter: RawFallbackFilter{
				GeoIP:  true,
				IPCIDR: []string{},
			},
			DefaultNameserver: []string{
				"114.114.114.114",
				"8.8.8.8",
			},
		},
	}

	profileJSON := json.NewDecoder(strings.NewReader(profile))
	profileJSON.UseNumber()
	if err := profileJSON.Decode(&rawCfg); err != nil {
		return nil, err
	}
	cfgJSON := json.NewDecoder(strings.NewReader(cfg))
	cfgJSON.UseNumber()
	if err := cfgJSON.Decode(&rawCfg); err != nil {
		return nil, err
	}

	return rawCfg, nil
}

func ParseRawConfig(rawCfg *RawConfig) (*CC.Config, error) {
	config := &CC.Config{}

	config.Experimental = &CC.Experimental{}

	general, err := parseGeneral(rawCfg)
	if err != nil {
		return nil, err
	}
	config.General = general

	proxies, providers, err := parseProxies(rawCfg)
	if err != nil {
		return nil, err
	}
	config.Proxies = proxies
	config.Providers = providers

	rules, err := parseRules(rawCfg, proxies)
	if err != nil {
		return nil, err
	}
	config.Rules = rules

	hosts, err := parseHosts(rawCfg)
	if err != nil {
		return nil, err
	}
	config.Hosts = hosts

	dnsCfg, err := parseDNS(rawCfg.DNS, hosts)
	if err != nil {
		return nil, err
	}
	config.DNS = dnsCfg

	config.Users = parseAuthentication(rawCfg.Authentication)

	return config, nil
}

func parseGeneral(cfg *RawConfig) (*CC.General, error) {
	externalUI := cfg.ExternalUI

	// checkout externalUI exist
	if externalUI != "" {
		externalUI = C.Path.Resolve(externalUI)

		if _, err := os.Stat(externalUI); os.IsNotExist(err) {
			return nil, fmt.Errorf("external-ui: %s not exist", externalUI)
		}
	}

	return &CC.General{
		Inbound: CC.Inbound{
			Port:        cfg.Port,
			SocksPort:   cfg.SocksPort,
			RedirPort:   cfg.RedirPort,
			MixedPort:   cfg.MixedPort,
			AllowLan:    cfg.AllowLan,
			BindAddress: cfg.BindAddress,
		},
		Controller: CC.Controller{
			ExternalController: cfg.ExternalController,
			ExternalUI:         cfg.ExternalUI,
			Secret:             cfg.Secret,
		},
		Mode:      cfg.Mode,
		LogLevel:  cfg.LogLevel,
		IPv6:      cfg.IPv6,
		Interface: cfg.Interface,
	}, nil
}

func parseProxies(cfg *RawConfig) (proxies map[string]C.Proxy, providersMap map[string]provider.ProxyProvider, err error) {
	proxies = make(map[string]C.Proxy)
	providersMap = make(map[string]provider.ProxyProvider)
	proxyList := []string{}
	proxiesConfig := cfg.Proxy
	groupsConfig := cfg.ProxyGroup
	providersConfig := cfg.ProxyProvider

	proxies["DIRECT"] = outbound.NewProxy(outbound.NewDirect())
	proxies["REJECT"] = outbound.NewProxy(outbound.NewReject())
	proxyList = append(proxyList, "DIRECT", "REJECT")

	// parse proxy
	for idx, mapping := range proxiesConfig {
		proxy, err := outbound.ParseProxy(mapping)
		if err != nil {
			return nil, nil, fmt.Errorf("proxy %d: %w", idx, err)
		}

		if _, exist := proxies[proxy.Name()]; exist {
			return nil, nil, fmt.Errorf("proxy %s is the duplicate name", proxy.Name())
		}
		proxies[proxy.Name()] = proxy
		proxyList = append(proxyList, proxy.Name())
	}

	// keep the original order of ProxyGroups in config file
	for idx, mapping := range groupsConfig {
		groupName, existName := mapping["name"].(string)
		if !existName {
			return nil, nil, fmt.Errorf("proxy group %d: missing name", idx)
		}
		proxyList = append(proxyList, groupName)
	}

	// check if any loop exists and sort the ProxyGroups
	if err := proxyGroupsDagSort(groupsConfig); err != nil {
		return nil, nil, err
	}

	// parse and initial providers
	for name, mapping := range providersConfig {
		if name == provider.ReservedName {
			return nil, nil, fmt.Errorf("can not defined a provider called `%s`", provider.ReservedName)
		}

		pd, err := provider.ParseProxyProvider(name, mapping)
		if err != nil {
			return nil, nil, fmt.Errorf("parse proxy provider %s error: %w", name, err)
		}

		providersMap[name] = pd
	}

	for _, provider := range providersMap {
		log.Infoln("Start initial provider %s", provider.Name())
		if err := provider.Initial(); err != nil {
			return nil, nil, fmt.Errorf("initial proxy provider %s error: %w", provider.Name(), err)
		}
	}

	// parse proxy group
	for idx, mapping := range groupsConfig {
		group, err := outboundgroup.ParseProxyGroup(mapping, proxies, providersMap)
		if err != nil {
			return nil, nil, fmt.Errorf("proxy group[%d]: %w", idx, err)
		}

		groupName := group.Name()
		if _, exist := proxies[groupName]; exist {
			return nil, nil, fmt.Errorf("proxy group %s: the duplicate name", groupName)
		}

		proxies[groupName] = outbound.NewProxy(group)
	}

	// initial compatible provider
	for _, pd := range providersMap {
		if pd.VehicleType() != provider.Compatible {
			continue
		}

		log.Infoln("Start initial compatible provider %s", pd.Name())
		if err := pd.Initial(); err != nil {
			return nil, nil, err
		}
	}

	ps := []C.Proxy{}
	for _, v := range proxyList {
		ps = append(ps, proxies[v])
	}
	hc := provider.NewHealthCheck(ps, "", 0)
	pd, _ := provider.NewCompatibleProvider(provider.ReservedName, ps, hc)
	providersMap[provider.ReservedName] = pd

	global := outboundgroup.NewSelector("GLOBAL", []provider.ProxyProvider{pd})
	proxies["GLOBAL"] = outbound.NewProxy(global)
	return proxies, providersMap, nil
}

func parseRules(cfg *RawConfig, proxies map[string]C.Proxy) ([]C.Rule, error) {
	rules := []C.Rule{}
	rulesConfig := cfg.Rule

	// parse rules
	for idx, line := range rulesConfig {
		rule := trimArr(strings.Split(line, ","))
		var (
			payload string
			target  string
			params  = []string{}
		)

		switch l := len(rule); {
		case l == 2:
			target = rule[1]
		case l == 3:
			payload = rule[1]
			target = rule[2]
		case l >= 4:
			payload = rule[1]
			target = rule[2]
			params = rule[3:]
		default:
			return nil, fmt.Errorf("rules[%d] [%s] error: format invalid", idx, line)
		}

		if _, ok := proxies[target]; !ok {
			return nil, fmt.Errorf("rules[%d] [%s] error: proxy [%s] not found", idx, line, target)
		}

		rule = trimArr(rule)
		params = trimArr(params)

		parsed, parseErr := R.ParseRule(rule[0], payload, target, params)
		if parseErr != nil {
			if parseErr == R.ErrPlatformNotSupport {
				log.Warnln("Rules[%d] [%s] don't support current OS, skip", idx, line)
				continue
			}
			return nil, fmt.Errorf("rules[%d] [%s] error: %s", idx, line, parseErr.Error())
		}

		rules = append(rules, parsed)
	}

	return rules, nil
}

func parseHosts(cfg *RawConfig) (*trie.DomainTrie, error) {
	tree := trie.New()

	// add default hosts
	if err := tree.Insert("localhost", net.IP{127, 0, 0, 1}); err != nil {
		log.Errorln("insert localhost to host error: %s", err.Error())
	}

	if len(cfg.Hosts) != 0 {
		for domain, ipStr := range cfg.Hosts {
			ip := net.ParseIP(ipStr)
			if ip == nil {
				return nil, fmt.Errorf("%s is not a valid IP", ipStr)
			}
			tree.Insert(domain, ip)
		}
	}

	return tree, nil
}

func hostWithDefaultPort(host string, defPort string) (string, error) {
	if !strings.Contains(host, ":") {
		host += ":"
	}

	hostname, port, err := net.SplitHostPort(host)
	if err != nil {
		return "", err
	}

	if port == "" {
		port = defPort
	}

	return net.JoinHostPort(hostname, port), nil
}

func parseNameServer(servers []string) ([]dns.NameServer, error) {
	nameservers := []dns.NameServer{}

	for idx, server := range servers {
		// parse without scheme .e.g 8.8.8.8:53
		if !strings.Contains(server, "://") {
			server = "udp://" + server
		}
		u, err := url.Parse(server)
		if err != nil {
			return nil, fmt.Errorf("DNS NameServer[%d] format error: %s", idx, err.Error())
		}

		var addr, dnsNetType string
		switch u.Scheme {
		case "udp":
			addr, err = hostWithDefaultPort(u.Host, "53")
			dnsNetType = "" // UDP
		case "tcp":
			addr, err = hostWithDefaultPort(u.Host, "53")
			dnsNetType = "tcp" // TCP
		case "tls":
			addr, err = hostWithDefaultPort(u.Host, "853")
			dnsNetType = "tcp-tls" // DNS over TLS
		case "https":
			clearURL := url.URL{Scheme: "https", Host: u.Host, Path: u.Path}
			addr = clearURL.String()
			dnsNetType = "https" // DNS over HTTPS
		default:
			return nil, fmt.Errorf("DNS NameServer[%d] unsupport scheme: %s", idx, u.Scheme)
		}

		if err != nil {
			return nil, fmt.Errorf("DNS NameServer[%d] format error: %s", idx, err.Error())
		}

		nameservers = append(
			nameservers,
			dns.NameServer{
				Net:  dnsNetType,
				Addr: addr,
			},
		)
	}
	return nameservers, nil
}

func parseFallbackIPCIDR(ips []string) ([]*net.IPNet, error) {
	ipNets := []*net.IPNet{}

	for idx, ip := range ips {
		_, ipnet, err := net.ParseCIDR(ip)
		if err != nil {
			return nil, fmt.Errorf("DNS FallbackIP[%d] format error: %s", idx, err.Error())
		}
		ipNets = append(ipNets, ipnet)
	}

	return ipNets, nil
}

func parseDNS(cfg RawDNS, hosts *trie.DomainTrie) (*CC.DNS, error) {
	if cfg.Enable && len(cfg.NameServer) == 0 {
		return nil, fmt.Errorf("if DNS configuration is turned on, NameServer cannot be empty")
	}

	dnsCfg := &CC.DNS{
		Enable:       cfg.Enable,
		Listen:       cfg.Listen,
		IPv6:         cfg.IPv6,
		EnhancedMode: cfg.EnhancedMode,
		FallbackFilter: CC.FallbackFilter{
			IPCIDR: []*net.IPNet{},
		},
	}
	var err error
	if dnsCfg.NameServer, err = parseNameServer(cfg.NameServer); err != nil {
		return nil, err
	}

	if dnsCfg.Fallback, err = parseNameServer(cfg.Fallback); err != nil {
		return nil, err
	}

	if len(cfg.DefaultNameserver) == 0 {
		return nil, errors.New("default nameserver should have at least one nameserver")
	}
	if dnsCfg.DefaultNameserver, err = parseNameServer(cfg.DefaultNameserver); err != nil {
		return nil, err
	}
	// check default nameserver is pure ip addr
	for _, ns := range dnsCfg.DefaultNameserver {
		host, _, err := net.SplitHostPort(ns.Addr)
		if err != nil || net.ParseIP(host) == nil {
			return nil, errors.New("default nameserver should be pure IP")
		}
	}

	if cfg.EnhancedMode == dns.FAKEIP {
		_, ipnet, err := net.ParseCIDR(cfg.FakeIPRange)
		if err != nil {
			return nil, err
		}

		var host *trie.DomainTrie
		// fake ip skip host filter
		if len(cfg.FakeIPFilter) != 0 {
			host = trie.New()
			for _, domain := range cfg.FakeIPFilter {
				host.Insert(domain, true)
			}
		}

		pool, err := fakeip.New(ipnet, 1000, host)
		if err != nil {
			return nil, err
		}

		dnsCfg.FakeIPRange = pool
	}

	dnsCfg.FallbackFilter.GeoIP = cfg.FallbackFilter.GeoIP
	if fallbackip, err := parseFallbackIPCIDR(cfg.FallbackFilter.IPCIDR); err == nil {
		dnsCfg.FallbackFilter.IPCIDR = fallbackip
	}
	dnsCfg.FallbackFilter.Domain = cfg.FallbackFilter.Domain

	if cfg.UseHosts {
		dnsCfg.Hosts = hosts
	}

	return dnsCfg, nil
}

func parseAuthentication(rawRecords []string) []auth.AuthUser {
	users := make([]auth.AuthUser, 0)
	for _, line := range rawRecords {
		userData := strings.SplitN(line, ":", 2)
		if len(userData) == 2 {
			users = append(users, auth.AuthUser{User: userData[0], Pass: userData[1]})
		}
	}
	return users
}

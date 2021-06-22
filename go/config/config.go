package config

import (
	"encoding/json"
	"github.com/jinzhu/copier"
	"net"
	"strings"

	"github.com/Dreamacro/clash/component/fakeip"
	"github.com/Dreamacro/clash/component/trie"
	clashConfig "github.com/Dreamacro/clash/config"
	"github.com/Dreamacro/clash/dns"
	"github.com/Dreamacro/clash/log"
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
	TProxyPort         int          `json:"tproxy-port"`
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
	Experimental  clashConfig.Experimental          `json:"experimental"`
	Proxy         []map[string]interface{}          `json:"proxies"`
	ProxyGroup    []map[string]interface{}          `json:"proxy-groups"`
	Rule          []string                          `json:"rules"`
}

// Parse config
func Parse(profile string, cfg string) (*clashConfig.Config, error) {
	rawCfg, err := UnmarshalRawConfig(profile, cfg)
	if err != nil {
		return nil, err
	}
	cfRawConfig := new(clashConfig.RawConfig)
	if err := copier.Copy(cfRawConfig, rawCfg); err != nil {
		return nil, err
	}
	return clashConfig.ParseRawConfig(cfRawConfig)
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

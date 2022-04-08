package config

import (
	"encoding/json"
	"github.com/jinzhu/copier"
	"net"
	"strings"

	"github.com/Dreamacro/clash/component/fakeip"
	"github.com/Dreamacro/clash/component/trie"
	clashConfig "github.com/Dreamacro/clash/config"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/dns"
	"github.com/Dreamacro/clash/log"
	T "github.com/Dreamacro/clash/tunnel"
)

// DNS config
type DNS struct {
	Enable            bool             `yaml:"enable" json:"enable"`
	IPv6              bool             `yaml:"ipv6" json:"ipv6"`
	NameServer        []dns.NameServer `yaml:"nameserver" json:"nameserver"`
	Fallback          []dns.NameServer `yaml:"fallback" json:"fallback"`
	FallbackFilter    FallbackFilter   `yaml:"fallback-filter" json:"fallback-filter"`
	Listen            string           `yaml:"listen" json:"listen"`
	EnhancedMode      C.DNSMode        `yaml:"enhanced-mode" json:"enhanced-mode"`
	DefaultNameserver []dns.NameServer `yaml:"default-nameserver" json:"default-nameserver"`
	FakeIPRange       *fakeip.Pool
	Hosts             *trie.DomainTrie
	NameServerPolicy  map[string]dns.NameServer
}

// FallbackFilter config
type FallbackFilter struct {
	GeoIP     bool         `yaml:"geoip" json:"geoip"`
	GeoIPCode string       `yaml:"geoip-code" json:"geoip-code"`
	IPCIDR    []*net.IPNet `yaml:"ipcidr" json:"ipcidr"`
	Domain    []string     `yaml:"domain" json:"domain"`
}

// Profile config
type Profile struct {
	StoreSelected bool `yaml:"store-selected" json:"store-selected"`
	StoreFakeIP   bool `yaml:"store-fake-ip" json:"store-fake-ip"`
}

type RawDNS struct {
	Enable            bool              `yaml:"enable" json:"enable"`
	IPv6              bool              `yaml:"ipv6" json:"ipv6"`
	UseHosts          bool              `yaml:"use-hosts" json:"use-hosts"`
	NameServer        []string          `yaml:"nameserver" json:"nameserver"`
	Fallback          []string          `yaml:"fallback" json:"fallback"`
	FallbackFilter    RawFallbackFilter `yaml:"fallback-filter" json:"fallback-filter"`
	Listen            string            `yaml:"listen" json:"listen"`
	EnhancedMode      C.DNSMode         `yaml:"enhanced-mode" json:"enhanced-mode"`
	FakeIPRange       string            `yaml:"fake-ip-range" json:"fake-ip-range"`
	FakeIPFilter      []string          `yaml:"fake-ip-filter" json:"fake-ip-filter"`
	DefaultNameserver []string          `yaml:"default-nameserver" json:"default-nameserver"`
	NameServerPolicy  map[string]string `yaml:"nameserver-policy" json:"nameserver-policy"`
}

type RawFallbackFilter struct {
	GeoIP     bool     `yaml:"geoip" json:"geoip"`
	GeoIPCode string   `yaml:"geoip-code" json:"geoip-code"`
	IPCIDR    []string `yaml:"ipcidr" json:"ipcidr"`
	Domain    []string `yaml:"domain" json:"domain"`
}

type RawConfig struct {
	Port               int          `yaml:"port" json:"port"`
	SocksPort          int          `yaml:"socks-port" json:"socks-port"`
	RedirPort          int          `yaml:"redir-port" json:"redir-port"`
	TProxyPort         int          `yaml:"tproxy-port" json:"tproxy-port"`
	MixedPort          int          `yaml:"mixed-port" json:"mixed-port"`
	Authentication     []string     `yaml:"authentication" json:"authentication"`
	AllowLan           bool         `yaml:"allow-lan" json:"allow-lan"`
	BindAddress        string       `yaml:"bind-address" json:"bind-address"`
	Mode               T.TunnelMode `yaml:"mode" json:"mode"`
	LogLevel           log.LogLevel `yaml:"log-level" json:"log-level"`
	IPv6               bool         `yaml:"ipv6" json:"ipv6"`
	ExternalController string       `yaml:"external-controller" json:"external-controller"`
	ExternalUI         string       `yaml:"external-ui" json:"external-ui"`
	Secret             string       `yaml:"secret" json:"secret"`
	Interface          string       `yaml:"interface-name" json:"interface-name"`
	RoutingMark        int          `yaml:"routing-mark" json:"routing-mark"`

	ProxyProvider map[string]map[string]any `yaml:"proxy-providers" json:"proxy-providers"`
	Hosts         map[string]string         `yaml:"hosts" json:"hosts"`
	DNS           RawDNS                    `yaml:"dns" json:"dns"`
	Experimental  clashConfig.Experimental  `yaml:"experimental" json:"experimental"`
	Profile       Profile                   `yaml:"profile" json:"profile"`
	Proxy         []map[string]any          `yaml:"proxies" json:"proxies"`
	ProxyGroup    []map[string]any          `yaml:"proxy-groups" json:"proxy-groups"`
	Rule          []string                  `yaml:"rules" json:"rules"`
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
		Proxy:          []map[string]any{},
		ProxyGroup:     []map[string]any{},
		DNS: RawDNS{
			Enable:      false,
			UseHosts:    true,
			FakeIPRange: "198.18.0.1/16",
			FallbackFilter: RawFallbackFilter{
				GeoIP:     true,
				GeoIPCode: "CN",
				IPCIDR:    []string{},
			},
			DefaultNameserver: []string{
				"114.114.114.114",
				"8.8.8.8",
			},
		},
		Profile: Profile{
			StoreSelected: true,
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

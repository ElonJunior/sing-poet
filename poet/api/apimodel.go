package api

import (
	"encoding/json"
	"regexp"
)

const (
	UserNotModified = "users not modified"
	NodeNotModified = "node not modified"
	RuleNotModified = "rules not modified"
)

// Config API config
type Config struct {
	APIHost string `mapstructure:"ApiHost"`
	NodeID  int    `mapstructure:"NodeID"`
	Key     string `mapstructure:"ApiKey"`
	InTag   string `mapstructure:"InTag"`

	NodeType            string  `mapstructure:"NodeType"`
	EnableVless         bool    `mapstructure:"EnableVless"`
	VlessFlow           string  `mapstructure:"VlessFlow"`
	Timeout             int     `mapstructure:"Timeout"`
	SpeedLimit          float64 `mapstructure:"SpeedLimit"`
	DeviceLimit         int     `mapstructure:"DeviceLimit"`
	RuleListPath        string  `mapstructure:"RuleListPath"`
	DisableCustomConfig bool    `mapstructure:"DisableCustomConfig"`
}

// NodeStatus Node status
type NodeStatus struct {
	CPU    float64
	Mem    float64
	Disk   float64
	Uptime uint64
}

// NodeType:
// - 与sing-box 大小写必须与struct定义相同
//
//   - Hysteria2 Clash.Meta最后release版本暂不支持obfs
//   - TUIC 小火箭对v5支持比较好的是从2005版开始
type NodeInfo struct {
	NodeType          string // Must be TUIC, Hysteria2, V2ray|VMess, Trojan, and Shadowsocks
	NodeID            int
	Port              uint32
	SpeedLimit        uint64 // Bps
	AlterID           uint16
	TransportProtocol string
	FakeType          string
	Host              string
	Path              string
	EnableTLS         bool
	EnableVless       bool
	VlessFlow         string
	CypherMethod      string
	ServerKey         string
	ServiceName       string
	Header            json.RawMessage
	// NameServerConfig  []*conf.NameServerConfig
	TrafficRate   float64 //流量倍率
	EnableREALITY bool
	REALITYConfig *REALITYConfig
}

type UserInfo struct {
	UID         int
	Email       string
	UUID        string
	Passwd      string
	Port        uint32
	AlterID     uint16
	Flow        string
	Method      string
	SpeedLimit  uint64 // Bps
	DeviceLimit int
}

type OnlineUser struct {
	UID int
	IP  string
}

type UserTraffic struct {
	UID      int
	Email    string
	Upload   int64
	Download int64
}

type ClientInfo struct {
	APIHost  string
	NodeID   int
	Key      string
	NodeType string
}

type DetectRule struct {
	ID      int
	Pattern *regexp.Regexp
}

type DetectResult struct {
	UID    int
	RuleID int
}

type REALITYConfig struct {
	Dest             string
	ProxyProtocolVer uint64
	ServerNames      []string
	PrivateKey       string
	MinClientVer     string
	MaxClientVer     string
	MaxTimeDiff      uint64
	ShortIds         []string
}

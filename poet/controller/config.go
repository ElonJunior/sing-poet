package controller

import "github.com/sagernet/sing-box/log"

type Config struct {
	PanelType      string `mapstructure:"PanelType"`
	UpdatePeriodic int    `mapstructure:"UpdatePeriodic"`

	ListenIP string `mapstructure:"ListenIP"`
	SendIP   string `mapstructure:"SendIP"`

	Logger log.Logger
}

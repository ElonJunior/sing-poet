package panel

import (
	"github.com/sagernet/sing-box/poet/api"
	"github.com/sagernet/sing-box/poet/controller"
)

type Config struct {
	Name        string         `mapstructure:"Name"`
	NodesConfig []*NodesConfig `mapstructure:"Nodes"`
}

type NodesConfig struct {
	// Panel type: SSpanel, NewV2board, PMpanel, Proxypanel, V2RaySocks, GoV2Panel
	PanelType        string             `mapstructure:"PanelType"`
	InTag            string             `mapstructure:"InTag"`
	OutTag           string             `mapstructure:"OutTag"`
	ApiConfig        *api.Config        `mapstructure:"ApiConfig"`
	ControllerConfig *controller.Config `mapstructure:"ControllerConfig"`
}

func getDefaultControllerConfig() *controller.Config {
	return &controller.Config{
		ListenIP:       "0.0.0.0",
		SendIP:         "0.0.0.0",
		UpdatePeriodic: 60,
	}
}

package panel

import (
	"context"
	"os"
	"sync"

	"dario.cat/mergo"
	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/poet/api"
	"github.com/sagernet/sing-box/poet/api/sspanel"
	"github.com/sagernet/sing-box/poet/controller"

	SS "github.com/sagernet/sing-box/poet/shortcuts"
)

// Panel Structure
type Panel struct {
	access      sync.Mutex
	panelConfig *Config
	Running     bool
	//Service
	Nodes []controller.Controller
}

func NewPanel(panelConfig *Config) *Panel {
	p := &Panel{
		panelConfig: panelConfig,
		Running:     false,
	}
	return p
}

// Start the panel
func (p *Panel) Start() {

	p.access.Lock()
	defer p.access.Unlock()

	//暂存对象
	ss := SS.Singleton()

	ctx := context.Background()

	log := ss.Logger
	log.Info("\n\tStart the panel..")

	//node counter
	nodeCounter := 0

	// Load Nodes config
	for _, nodeConfig := range p.panelConfig.NodesConfig {
		log.Info("processing NodeId: ", nodeConfig.ApiConfig.NodeID)

		//inbound
		var inbound *adapter.Inbound

		inTag := nodeConfig.InTag
		outTag := nodeConfig.OutTag
		if in, ok := ss.Inbounds[inTag]; ok {
			inbound = in
		} else {
			log.Warn("unkown inbound with Tag:", nodeConfig.ApiConfig.InTag)
			continue
		}

		//panel webapi
		var apiClient api.API

		switch nodeConfig.PanelType {
		case "SSpanel":
			apiClient = sspanel.New(nodeConfig.ApiConfig)
		default:
			log.Panic("unsupport panel type: ", nodeConfig.PanelType)
			continue
		}

		// Register controller service
		controllerConfig := getDefaultControllerConfig()
		controllerConfig.PanelType = nodeConfig.PanelType

		if nodeConfig.ControllerConfig != nil {
			if err := mergo.Merge(controllerConfig, nodeConfig.ControllerConfig, mergo.WithOverride); err != nil {
				log.Panic("read Controller Config Failed!")
				continue
			}
		}
		controllerService := controller.New(ctx, *controllerConfig, inbound, apiClient, log)
		p.Nodes = append(p.Nodes, *controllerService)

		//set Contrl
		ss.SetContrl(controllerService, inTag, outTag)

		nodeCounter++
	}

	//empty node
	if nodeCounter == 0 {
		log.Panic("no Node with InTag found!")
		os.Exit(1)
	}

	// Start all the nodes
	for _, node := range p.Nodes {
		err := node.Start()
		if err != nil {
			log.Panic("Panel Start fialed: ", err)
			os.Exit(1)
		}
	}

	p.Running = true
}

// Close the panel
func (p *Panel) Close() {
	p.access.Lock()
	defer p.access.Unlock()

	for _, node := range p.Nodes {
		node.Close()
	}

	p.Running = false
}

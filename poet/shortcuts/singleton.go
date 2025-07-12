package shortcuts

import (
	"sync"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/poet/controller"
)

type singleton struct {
	globalMap sync.Map

	Logger log.Logger

	//inbound
	Inbounds map[string]*adapter.Inbound
	//inTag and outTag map
	InContrls  map[string]*controller.Controller
	OutContrls map[string]*controller.Controller
}

var instance *singleton
var once sync.Once

func Singleton() *singleton {
	once.Do(func() {
		instance = &singleton{
			Inbounds:   make(map[string]*adapter.Inbound),
			InContrls:  make(map[string]*controller.Controller),
			OutContrls: make(map[string]*controller.Controller),
		}
	})
	return instance
}

func GetObject(k string) (interface{}, bool) {
	return Singleton().globalMap.Load(k)
}

func SetObject(k string, v interface{}) {
	Singleton().globalMap.Store(k, v)
}

func Log() log.Logger {
	return Singleton().Logger
}

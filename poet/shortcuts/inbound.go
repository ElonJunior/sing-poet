package shortcuts

import (
	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/poet/controller"
)

func (ss *singleton) SetInboud(in *adapter.Inbound, tag string) {
	// ins := *in
	// tags := ins.Tag()
	ss.Inbounds[tag] = in
}

func (ss *singleton) GetInoud(tag string) *adapter.Inbound {
	if in, ok := ss.Inbounds[tag]; ok {
		return in
	}
	return nil
}

func (ss *singleton) SetContrl(ct *controller.Controller, inTag string, outTag string) {
	ss.InContrls[inTag] = ct
	ss.OutContrls[outTag] = ct
}

func (ss *singleton) GetContrlWithInTag(inTag string) *controller.Controller {
	if ct, ok := ss.InContrls[inTag]; ok {
		return ct
	}
	return nil
}
func (ss *singleton) GetContrlWithOutTag(outTag string) *controller.Controller {
	if ct, ok := ss.OutContrls[outTag]; ok {
		return ct
	}
	return nil
}

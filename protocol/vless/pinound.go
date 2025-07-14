package vless

import (
	"fmt"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-box/poet/api"
)

// 确保实现 UserRefresher 接口
var _ adapter.PInbound = (*Inbound)(nil)

// RefreshUsers 更新用户列表
func (h *Inbound) RefreshUsers(users *[]api.UserInfo, nodeInfo *api.NodeInfo) error {

	// 1. 准备更新参数
	ulen := len(*users)
	indices := make([]int, ulen)
	uuids := make([]string, ulen)
	flows := make([]string, ulen)
	opUsers := make([]option.VLESSUser, ulen)

	for i, user := range *users {
		indices[i] = i
		uuids[i] = user.UUID

		flow := user.Flow
		if flow == "" {
			// 如果用户没有设置Flow，使用节点配置的Flow
			flow = nodeInfo.VlessFlow
		}
		flows[i] = flow

		opUsers[i] = option.VLESSUser{
			Name: fmt.Sprintf("u%d", user.UID),
			UUID: user.UUID,
			Flow: flow,
		}
	}

	// 2. 更新服务
	h.service.UpdateUsers(indices, uuids, flows)

	// 3. 更新本地用户列表
	h.users = opUsers

	return nil
}

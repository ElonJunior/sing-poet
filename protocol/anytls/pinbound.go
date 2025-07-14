package anytls

import (
	"fmt"

	anytls "github.com/anytls/sing-anytls"
	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/poet/api"
)

// 确保实现 UserRefresher 接口
var _ adapter.PInbound = (*Inbound)(nil)

// RefreshUsers 更新用户列表
func (h *Inbound) RefreshUsers(users *[]api.UserInfo, nodeInfo *api.NodeInfo) error {
	// 1. 准备更新参数
	ulen := len(*users)
	opUsers := make([]anytls.User, ulen)

	for i, user := range *users {
		opUsers[i] = anytls.User{
			Name:     fmt.Sprintf("u%d", user.UID),
			Password: user.Passwd,
		}
	}

	// 2. 更新服务用户
	h.service.UpdateUsers(opUsers)

	return nil
}

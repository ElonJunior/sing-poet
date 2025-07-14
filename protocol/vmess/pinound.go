package vmess

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
	userList := make([]int, ulen)
	userIdList := make([]string, ulen)
	alterIdList := make([]int, ulen)
	opUsers := make([]option.VMessUser, ulen)

	for i, user := range *users {
		userList[i] = i
		userIdList[i] = user.UUID

		alterId := int(user.AlterID)
		if alterId == 0 {
			alterId = int(nodeInfo.AlterID)
		}
		alterIdList[i] = alterId

		opUsers[i] = option.VMessUser{
			Name:    fmt.Sprintf("u%d", user.UID),
			UUID:    user.UUID,
			AlterId: alterId,
		}

	}

	// 2. 更新本地用户列表
	err := h.service.UpdateUsers(userList, userIdList, alterIdList)
	h.users = opUsers

	return err
}

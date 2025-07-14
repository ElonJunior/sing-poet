package trojan

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

	userList := make([]int, len(*users))
	passwordList := make([]string, len(*users))
	opUsers := make([]option.TrojanUser, len(*users))

	for i, user := range *users {
		userList[i] = i
		passwordList[i] = user.Passwd

		opUsers[i] = option.TrojanUser{
			Name:     fmt.Sprintf("u%d", user.UID),
			Password: user.Passwd,
		}
	}

	err := h.service.UpdateUsers(userList, passwordList)
	if err != nil {
		return err
	}

	h.users = opUsers

	return nil
}

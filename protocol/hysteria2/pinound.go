package hysteria2

import (
	"fmt"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/poet/api"
)

// 确保实现 UserRefresher 接口
var _ adapter.PInbound = (*Inbound)(nil)

// RefreshUsers 更新用户列表
func (h *Inbound) RefreshUsers(users *[]api.UserInfo, nodeInfo *api.NodeInfo) error {

	// 1. 准备更新参数
	ulen := len(*users)
	userList := make([]int, ulen)
	userPasswordList := make([]string, ulen)
	userNameList := make([]string, ulen)

	for i, user := range *users {
		userNameList[i] = fmt.Sprintf("u%d", user.UID)

		userList[i] = i
		userPasswordList[i] = user.Passwd
	}

	// 2. 更新服务
	h.service.UpdateUsers(userList, userPasswordList)
	// 3. 更新本地用户列表
	h.userNameList = userNameList

	return nil
}

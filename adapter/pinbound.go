package adapter

import "github.com/sagernet/sing-box/poet/api"

// UserRefresher 接口定义用户刷新功能
// 这是一个可选接口，不是所有 Inbound 都需要实现
type PInbound interface {
	// RefreshUsers 更新用户数据
	RefreshUsers(users *[]api.UserInfo, nodeInfo *api.NodeInfo) error
}

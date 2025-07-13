package adapter

// UserRefresher 接口定义用户刷新功能
// 这是一个可选接口，不是所有 Inbound 都需要实现
type PInbound interface {
	// RefreshUsers 更新用户数据
	// users 参数类型取决于具体实现（如 []option.VMessUser, []option.ShadowsocksUser）
	RefreshUsers(users any) error
}

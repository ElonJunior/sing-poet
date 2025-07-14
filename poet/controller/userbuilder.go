package controller

import (
	"fmt"
	"strings"

	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-box/poet/api"
)

// 修改 UserBuilder 类型定义，接收 *api.NodeInfo 参数
type UserBuilder func(*api.NodeInfo, *[]api.UserInfo) any

// 创建用户构建器映射表
func (c *Controller) createUserBuilders() map[string]func(*api.NodeInfo, *[]api.UserInfo) any {
	return map[string]func(*api.NodeInfo, *[]api.UserInfo) any{
		"vmess": func(nodeInfo *api.NodeInfo, users *[]api.UserInfo) any {
			return c.buildVMessUser(nodeInfo, users)
		},
		"vless": func(nodeInfo *api.NodeInfo, users *[]api.UserInfo) any {
			return c.buildVlessUser(nodeInfo, users)
		},
		"tuic": func(nodeInfo *api.NodeInfo, users *[]api.UserInfo) any {
			return c.buildTUICUser(users)
		},
		"hysteria2": func(nodeInfo *api.NodeInfo, users *[]api.UserInfo) any {
			return c.buildHysteria2User(users)
		},
		"trojan": func(nodeInfo *api.NodeInfo, users *[]api.UserInfo) any {
			return c.buildTrojanUser(users)
		},
		"anytls": func(nodeInfo *api.NodeInfo, users *[]api.UserInfo) any {
			return c.buildAnyTLSUser(users)
		},
		// 添加更多协议支持...
	}
}

// 修改 BuildUsers 方法签名和实现
func (c *Controller) BuildUsers(nodeInfo *api.NodeInfo, userInfo *[]api.UserInfo) (any, error) {
	// 获取构建器映射
	builders := c.createUserBuilders()

	// 获取对应的构建器
	builder, exist := builders[strings.ToLower(nodeInfo.NodeType)]
	if !exist {
		return nil, fmt.Errorf("unsupported node type: %s", nodeInfo.NodeType)
	}
	return builder(nodeInfo, userInfo), nil
}

func (c *Controller) buildUserHash(user *api.UserInfo) string {
	var hash string
	switch c.panelType {
	case "SSpanel":
		hash = fmt.Sprintf("u%d", user.UID)
	default:
		hash = user.UUID
	}

	return hash
}

func (c *Controller) buildTUICUser(userInfo *[]api.UserInfo) (users []*option.TUICUser) {
	users = make([]*option.TUICUser, len(*userInfo))
	for i, user := range *userInfo {
		//debug
		// c.log(fmt.Sprintf("Name:%s\tUUID:%s\tPwd:%s", fmt.Sprintf("u%d", user.UID), user.UUID, user.Passwd), "debug")

		users[i] = &option.TUICUser{
			UUID:     user.UUID,
			Name:     c.buildUserHash(&user),
			Password: user.Passwd,
		}
	}
	return users
}

func (c *Controller) buildHysteria2User(userInfo *[]api.UserInfo) (users []*option.Hysteria2User) {
	users = make([]*option.Hysteria2User, len(*userInfo))
	for i, user := range *userInfo {
		users[i] = &option.Hysteria2User{
			Name:     c.buildUserHash(&user),
			Password: user.Passwd,
		}
	}
	return users
}

func (c *Controller) buildVMessUser(nodeInfo *api.NodeInfo, userInfo *[]api.UserInfo) (users []option.VMessUser) {
	users = make([]option.VMessUser, len(*userInfo))
	for i, user := range *userInfo {
		alterId := int(user.AlterID)
		if alterId == 0 {
			alterId = int(nodeInfo.AlterID)
		}

		users[i] = option.VMessUser{
			Name:    c.buildUserHash(&user),
			UUID:    user.UUID,
			AlterId: alterId,
		}
	}
	return users
}

func (c *Controller) buildVlessUser(nodeInfo *api.NodeInfo, userInfo *[]api.UserInfo) (users []option.VLESSUser) {
	users = make([]option.VLESSUser, len(*userInfo))
	for i, user := range *userInfo {
		flow := user.Flow
		if flow == "" {
			// 如果用户没有设置Flow，使用节点配置的Flow
			flow = nodeInfo.VlessFlow
		}

		users[i] = option.VLESSUser{
			Name: c.buildUserHash(&user),
			UUID: user.UUID,
			Flow: flow,
		}
	}
	return users
}

func (c *Controller) buildTrojanUser(userInfo *[]api.UserInfo) (users []*option.TrojanUser) {
	users = make([]*option.TrojanUser, len(*userInfo))
	for i, user := range *userInfo {
		users[i] = &option.TrojanUser{
			Name:     c.buildUserHash(&user),
			Password: user.Passwd,
		}
	}
	return users
}

func (c *Controller) buildAnyTLSUser(userInfo *[]api.UserInfo) (users []*option.AnyTLSUser) {
	users = make([]*option.AnyTLSUser, len(*userInfo))
	for i, user := range *userInfo {
		users[i] = &option.AnyTLSUser{
			Name:     c.buildUserHash(&user),
			Password: user.Passwd,
		}
	}
	return users
}

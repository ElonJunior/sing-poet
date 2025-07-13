package controller

import (
	"fmt"
	"strings"

	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-box/poet/api"
)

// 定义用户构建器函数类型
type UserBuilder func([]api.UserInfo) any

// 创建用户构建器映射表
func (c *Controller) createUserBuilders() map[string]func(*[]api.UserInfo) any {
	return map[string]func(*[]api.UserInfo) any{
		"vmess": func(users *[]api.UserInfo) any {
			return c.buildVMessUser(users)
		},
		"tuic": func(users *[]api.UserInfo) any {
			return c.buildTUICUser(users)
		},
		"hysteria2": func(users *[]api.UserInfo) any {
			return c.buildHysteria2User(users)
		},
		// 添加更多协议支持...
	}
}

// 统一构建用户数据
func (c *Controller) BuildUsers(userType string, userInfo *[]api.UserInfo) (any, error) {
	// 获取构建器映射
	builders := c.createUserBuilders()

	// 获取对应的构建器
	builder, exists := builders[strings.ToLower(userType)]
	if !exists {
		return nil, fmt.Errorf("unsupported user type: %s", userType)
	}

	// 执行构建
	return builder(userInfo), nil
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

func (c *Controller) buildVMessUser(userInfo *[]api.UserInfo) (users []option.VMessUser) {
	users = make([]option.VMessUser, len(*userInfo))
	for i, user := range *userInfo {
		users[i] = option.VMessUser{
			Name:    c.buildUserHash(&user),
			UUID:    user.UUID,
			AlterId: int(user.AlterID),
		}
	}
	return users
}

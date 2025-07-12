package controller

import (
	"fmt"

	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-box/poet/api"
)

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

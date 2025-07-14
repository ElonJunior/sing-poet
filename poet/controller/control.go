package controller

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/poet/api"
)

func (c *Controller) syncUserList() error {
	//get list
	userInfo, err := c.apiClient.GetUserList()
	if err != nil {
		//api.UserNotModified
		if err.Error() == api.UserNotModified {
			c.log("GetUserList: 304 UserNotModified", "info")
			return nil
		}

		return err
	}

	//获得 更新两张 用户map表
	added, deleted := c.compareUsersMapReDeleted(&c.usersMap, userInfo)
	adLen, deLen := len(added), len(deleted)
	c.log(fmt.Sprintf("Sync Users added: %d deleted: %d", adLen, deLen), "info")
	if adLen == 0 && deLen == 0 {
		return nil
	}

	if deLen > 0 {
		for hash := range deleted {
			c.log(fmt.Sprintf("DeleteUser: %s", hash), "debug")
			err := c.author.DelUser(hash)
			if err != nil {
				c.log(fmt.Sprintf("failed deleteUser: %s %v", hash, err.Error()), "error")
				//return err
			}
		}
	}
	if adLen > 0 {
		var uid int
		for hash := range added {
			if user, ok := c.usersMap[hash]; ok {
				uid = user.UID
			} else {
				c.log(fmt.Sprintf("usersMap none hash: %s", hash), "error")
				continue
			}

			// c.log(fmt.Sprintf("AddUser: %s UID: %d", hash, uid), "debug")
			err := c.author.AddUser(hash, uid)
			if err != nil {
				c.log(fmt.Sprintf("failed addUser: %s %v", hash, err.Error()), "error")
				//return err
			}
		}
	}

	// update inbound server userList
	// 1. 检查是否实现 UserRefresher 接口
	in := *c.inbound
	refresher, ok := in.(adapter.PInbound)
	if !ok {
		c.log(fmt.Sprintf("unsupported node type: %s", c.nodeInfo.NodeType), "error")
		return errors.New("inbound type does not support user refresh")
	}

	// // 2. 构建适合该类型的用户数据
	// users, err := c.BuildUsers(c.nodeInfo, userInfo)
	// if err != nil {
	// 	c.log(fmt.Sprintf("Failed to build users for node type %s, error: %v", c.nodeInfo.NodeType, err), "error")
	// 	return err
	// }

	// 3. 刷新用户
	if err := refresher.RefreshUsers(userInfo, c.nodeInfo); err != nil {
		c.log(fmt.Sprintf("Failed to refresh users for node type %s, error: %v", c.nodeInfo.NodeType, err), "error")
		return err
	}

	c.log(fmt.Sprintf("final UsersMap: %d \tuserInfo: %d ", len(c.usersMap), len(*userInfo)), "info")
	return nil
}

// 记录流量
func (c *Controller) LoadOrCreateUserCounter(ctx context.Context, metadata adapter.InboundContext) (send, recv *atomic.Int64, err error) {
	//每次 goroutine 会预先生成用户Map 以及 author.users

	hash := metadata.User
	meter, found := c.author.users.Load(hash)
	if !found {
		return nil, nil, fmt.Errorf("hash:%s is not in user map", hash)
	}

	user := meter.(*User)
	//add IP
	user.AddIP(metadata.Source.AddrString())

	//pointer
	sendPtr, recvPtr := user.GetTrafficPointer()
	return sendPtr, recvPtr, nil
}

// 返回需要删除的 hash map
func (c *Controller) compareUsersMapReDeleted(old *map[string]*api.UserInfo, new *[]api.UserInfo) (added, deleted map[string]int) {
	olen := len(*old)
	nlen := len(*new)
	if olen == 0 && nlen == 0 {
		return nil, nil
	}

	// fmt.Printf("compare olen:%d nlen:%d\n", olen, nlen)

	deleted = make(map[string]int)
	oldMap := *old
	for hash := range oldMap {
		deleted[hash] = 1
	}

	//新用户列表 空, 删除全部
	if nlen == 0 {
		// fmt.Printf("compare nlen=0 deleted:%d\n", len(deleted))
		return nil, deleted
	}

	//循环新用户列表
	added = make(map[string]int)
	for _, user := range *new {
		// 1. slice元素是数组元素引用,取址后都是同一块内存
		// 2. for range默认取地址,多个迭代变量指向同一内存
		// 3. 应该在遍历时进行值拷贝,使每个迭代变量隔离
		// Go语言中的切片遍历需要注意这一点,以避免修改都作用在同一内存上的问题

		hash := c.buildUserHash(&user)
		// fmt.Printf("compare user hash:%s\n", hash)

		//老map为空 added+1 deleted=0
		if olen == 0 {
			added[hash] = 1
			u := user //赋值之前 进行值拷贝
			oldMap[hash] = &u
			// fmt.Printf("compare olen=0 added:%d\n", len(added))
			continue
		}

		//判断map
		if _, ok := oldMap[hash]; ok {
			// 存在key
			delete(deleted, hash)
			// fmt.Printf("compare hash exists deleted: %s\n", hash)
		} else {
			// 不存在key
			added[hash] = 1
			u := user //赋值之前 进行值拷贝
			oldMap[hash] = &u
			// fmt.Printf("compare added hash: %s\n", hash)
		}
	}

	//删除 多余的kv
	if len(deleted) > 0 {
		// fmt.Printf("compare len(deleted)>0 deleted:%d\n", len(deleted))
		for hash := range deleted {
			delete(oldMap, hash)
			// fmt.Printf("compare deleted hash: %s\n", hash)
		}
	}

	return added, deleted
}

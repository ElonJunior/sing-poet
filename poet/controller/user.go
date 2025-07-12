package controller

//https://github.com/bobby4k/trojan-go/blob/master/statistic/memory/memory.go
//atomic

import (
	"context"
	"fmt"
	"sync"

	// "sync/atomic"
	"github.com/sagernet/sing/common/atomic"

	"golang.org/x/time/rate"
)

const Name = "MEMORY"

type User struct {
	// WARNING: do not change the order of these fields.
	// 64-bit fields that use `sync/atomic` package functions
	// must be 64-bit aligned on 32-bit systems.
	// Reference: https://github.com/golang/go/issues/599
	// Solution: https://github.com/golang/go/issues/11891#issuecomment-433623786
	sent      atomic.Int64
	recv      atomic.Int64
	lastSent  atomic.Int64
	lastRecv  atomic.Int64
	sendSpeed atomic.Int64
	recvSpeed atomic.Int64

	hash string
	UID  int

	ipTable     sync.Map
	ipNum       atomic.Int32
	maxIPNum    int
	limiterLock sync.RWMutex
	sendLimiter *rate.Limiter
	recvLimiter *rate.Limiter

	ctx    context.Context
	cancel context.CancelFunc
}

func (u *User) Close() error {
	u.ResetTraffic()
	u.cancel()
	return nil
}

func (u *User) AddIP(ip string) bool {
	//先加入ip
	_, found := u.ipTable.Load(ip)
	if found {
		return true //无需限制, 再table中
	}

	//判断是否受限
	if u.maxIPNum > 0 {
		ipn := u.ipNum.Load()
		if int(ipn)+1 > u.maxIPNum {
			return false //超过ip限制
		}
	}

	u.ipTable.Store(ip, true)
	u.ipNum.Add(1)

	return true
}

func (u *User) DelIP(ip string) bool {
	//先载入
	_, found := u.ipTable.Load(ip)
	if !found {
		return false
	}

	// if u.maxIPNum <= 0 {
	// 	return true
	// }

	u.ipTable.Delete(ip)
	u.ipNum.Add(-1)
	return true
}

func (u *User) GetIPNumber() int {
	return int(u.ipNum.Load())
}

func (u *User) SetIPLimit(n int) {
	u.maxIPNum = n
}

func (u *User) GetIPLimit() int {
	return u.maxIPNum
}

func (u *User) AddTraffic(sent, recv int64) {
	u.limiterLock.RLock()
	defer u.limiterLock.RUnlock()

	if u.sendLimiter != nil && sent >= 0 {
		u.sendLimiter.WaitN(u.ctx, int(sent))
	} else if u.recvLimiter != nil && recv >= 0 {
		u.recvLimiter.WaitN(u.ctx, int(recv))
	}

	u.sent.Add(sent)
	u.recv.Add(recv)
	// atomic.AddUint64(&u.sent, uint64(sent))
	// atomic.AddUint64(&u.recv, uint64(recv))
}

func (u *User) SetSpeedLimit(send, recv int64) {
	u.limiterLock.Lock()
	defer u.limiterLock.Unlock()

	if send <= 0 {
		u.sendLimiter = nil
	} else {
		u.sendLimiter = rate.NewLimiter(rate.Limit(send), int(send)*2)
	}
	if recv <= 0 {
		u.recvLimiter = nil
	} else {
		u.recvLimiter = rate.NewLimiter(rate.Limit(recv), int(recv)*2)
	}
}

func (u *User) GetSpeedLimit() (send, recv int) {
	u.limiterLock.RLock()
	defer u.limiterLock.RUnlock()

	if u.sendLimiter != nil {
		send = int(u.sendLimiter.Limit())
	}
	if u.recvLimiter != nil {
		recv = int(u.recvLimiter.Limit())
	}
	return
}

func (u *User) Hash() string {
	return u.hash
}

func (u *User) SetTraffic(send, recv int64) {
	u.sent.Store(send)
	u.recv.Store(recv)
	// atomic.StoreUint64(&u.sent, send)
	// atomic.StoreUint64(&u.recv, recv)
}

func (u *User) GetTraffic() (int64, int64) {
	return u.sent.Load(), u.recv.Load()
}

func (u *User) GetTrafficPointer() (send, recv *atomic.Int64) {
	return &u.sent, &u.recv
}

func (u *User) ResetTraffic() (int64, int64) {
	sent := u.sent.Swap(0)
	recv := u.recv.Swap(0)
	u.lastSent.Store(0)
	u.lastRecv.Store(0)
	return sent, recv
}

func (u *User) ResetIPTable() []string {
	var iplist []string

	u.ipTable.Range(func(key, value interface{}) bool {
		if ip, ok := key.(string); ok {
			iplist = append(iplist, ip)

			//reset
			u.ipTable.Delete(ip)
			u.ipNum.Add(-1)
		}
		return true
	})

	return iplist
}

// // 对应下面携程模式的快速修改
// func (u *User) TardyAddTraffic(send, recv int64) {
// 	if recv > 0 {
// 		atomic.AddUint64(&u.recv, recv)
// 	}
// 	if send > 0 {
// 		atomic.AddUint64(&u.sent, send)
// 	}
// }

// func (u *User) speedUpdater() {
// 	ticker := time.NewTicker(time.Second)
// 	for {
// 		select {
// 		case <-u.ctx.Done():
// 			return
// 		case <-ticker.C:
// 			sent, recv := u.GetTraffic()
// 			atomic.StoreUint64(&u.sendSpeed, sent-u.lastSent)
// 			atomic.StoreUint64(&u.recvSpeed, recv-u.lastRecv)
// 			atomic.StoreUint64(&u.lastSent, sent)
// 			atomic.StoreUint64(&u.lastRecv, recv)
// 		}
// 	}
// }

func (u *User) GetSpeed() (int64, int64) {
	return u.sendSpeed.Load(), u.recvSpeed.Load()
}

type Authenticator struct {
	users sync.Map
	ctx   context.Context
}

// func (a *Authenticator) AuthUser(hash string) (bool, statistic.User) {
// 	if user, found := a.users.Load(hash); found {
// 		return true, user.(*User)
// 	}
// 	return false, nil
// }

func (a *Authenticator) AddUser(hash string, uid int) error {
	if _, found := a.users.Load(hash); found {
		return fmt.Errorf("hash " + hash + " is already exist")
	}
	ctx, cancel := context.WithCancel(a.ctx)
	meter := &User{
		hash:   hash,
		UID:    uid,
		ctx:    ctx,
		cancel: cancel,
	}
	// go meter.speedUpdater()
	a.users.Store(hash, meter)
	return nil
}

func (a *Authenticator) DelUser(hash string) error {
	meter, found := a.users.Load(hash)
	if !found {
		return fmt.Errorf("hash " + hash + " not found")
	}
	meter.(*User).Close()
	a.users.Delete(hash)
	return nil
}

// func (a *Authenticator) AddTraffic(hash string, send, recv uint64) error {
// 	meter, found := a.users.Load(hash)
// 	if !found {
// 		return fmt.Errorf("hash " + hash + " not found")
// 	}
// 	meter.(*User).TardyAddTraffic(send, recv)
// 	return nil
// }

// 返回指针列表 复制sync.Map内部的互斥锁,可能会导致并发问题
func (a *Authenticator) ListUsers() []*User {
	result := make([]*User, 0)

	a.users.Range(func(key, value interface{}) bool {
		if u, ok := value.(*User); ok {
			result = append(result, u)
		}
		return true
	})

	return result
}

func (a *Authenticator) Close() error {
	return nil
}

func NewAuthenticator(ctx context.Context) (*Authenticator, error) {
	u := &Authenticator{
		ctx: ctx,
	}
	return u, nil
}

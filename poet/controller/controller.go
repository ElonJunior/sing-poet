package controller

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/poet/api"
	"github.com/sagernet/sing-box/poet/common/serverstatus"
	"github.com/sagernet/sing-box/poet/common/task"
)

type Controller struct {
	config     *Config
	clientInfo api.ClientInfo
	apiClient  api.API
	nodeInfo   *api.NodeInfo
	Tag        string
	// userList   *[]api.UserInfo
	usersMap map[string]*api.UserInfo

	tasks []periodicTask
	// limitedUsers map[api.UserInfo]LimitInfo
	// warnedUsers  map[api.UserInfo]int
	// Panel type: SSpanel, NewV2board, PMpanel, Proxypanel, V2RaySocks, GoV2Panel
	panelType string
	// stm          stats.Manager
	// dispatcher   *mydispatcher.DefaultDispatcher
	startAt time.Time
	inbound *adapter.Inbound
	logger  log.Logger

	//用户统计与验证
	author *Authenticator
	ctx    context.Context

	// environment = development testing staging production
	env string
}

type periodicTask struct {
	tag string
	*task.Periodic
}

// New return a Controller service with default parameters.
func New(ctx context.Context, config Config, inbound *adapter.Inbound, apiClient api.API, logger log.Logger) *Controller {

	controller := &Controller{
		panelType: config.PanelType,
		config:    &config,
		apiClient: apiClient,
		// userList:  &[]api.UserInfo{},
		usersMap: make(map[string]*api.UserInfo),

		startAt: time.Now(),
		inbound: inbound,

		ctx:    ctx,
		logger: logger,
	}
	if author, err := NewAuthenticator(ctx); err == nil {
		controller.author = author
	} else {
		logger.Panic("err author init fail!")
	}

	return controller
}

// Start implement the Start() function of the service interface
func (c *Controller) Start() error {

	// log := c.config.Logger

	c.clientInfo = c.apiClient.Describe()

	// First fetch Node Info
	newNodeInfo, err := c.apiClient.GetNodeInfo()
	if err != nil {
		return err
	}
	if newNodeInfo.Port == 0 {
		return errors.New("server port must > 0")
	}

	//简单验证信息
	// if newNodeInfo.NodeType != c.clientInfo.NodeType {
	// 	return errors.New("node type must be same")
	// }

	c.nodeInfo = newNodeInfo
	c.Tag = c.buildNodeTag() //== inTag

	//sync controller user list
	err = c.syncUserList()
	if err != nil {
		return err
	}

	// Add periodic tasks
	c.tasks = append(c.tasks,
		periodicTask{
			tag: "node monitor",
			Periodic: &task.Periodic{
				Interval: time.Duration(c.config.UpdatePeriodic) * time.Second,
				Execute:  c.nodeInfoMonitor,
			}},
		periodicTask{
			tag: "user monitor",
			Periodic: &task.Periodic{
				Interval: time.Duration(c.config.UpdatePeriodic) * time.Second,
				Execute:  c.userInfoMonitor,
			}},
	)

	// Start periodic tasks
	for i := range c.tasks {
		c.log("start periodic task: "+c.tasks[i].tag, "info")
		go c.tasks[i].Start()
	}
	return nil
}

// Close implement the Close() function of the service interface
func (c *Controller) Close() error {
	log := c.config.Logger

	for i := range c.tasks {
		if c.tasks[i].Periodic != nil {
			if err := c.tasks[i].Periodic.Close(); err != nil {
				// log.Panic("%s %s periodic task close failed: %s", c.logPrefix(), c.tasks[i].tag, err)
				log.Panic("periodic task close failed: %s", err)
				return err
			}
		}
	}

	return nil
}

func (c *Controller) buildNodeTag() string {
	return fmt.Sprintf("%s_%s_%d", c.nodeInfo.NodeType, c.config.ListenIP, c.nodeInfo.Port)
}

func (c *Controller) Env() string {
	if c.env == "" {
		if os.Getenv("WSL_DISTRO_NAME") != "" {
			c.env = "development"
		} else {
			c.env = "staging"
		}
	}
	return c.env
}

func (c *Controller) log(message any, level string) {
	message = fmt.Sprintf("Node[%d] %v", c.nodeInfo.NodeID, message)

	switch level {
	case "debug":
		if c.Env() == "development" {
			fmt.Println(time.Now().Format("01-02 15:04:05 "), message) //std output
		}
		c.logger.Debug(message)
	case "warn":
		c.logger.Warn(message)
	case "error":
		c.logger.Error(message)
	case "panic":
		c.logger.Panic(message)
	default:
		c.logger.Info(message)
	}
}

func (c *Controller) nodeInfoMonitor() (err error) {
	// delay to start
	if time.Since(c.startAt) < time.Duration(c.config.UpdatePeriodic)*time.Second {
		return nil
	}

	// First fetch Node Info
	var nodeInfoChanged = true
	newNodeInfo, err := c.apiClient.GetNodeInfo()
	if err != nil {
		if err.Error() == api.NodeNotModified {
			nodeInfoChanged = false
			newNodeInfo = c.nodeInfo
		} else {
			c.log(fmt.Sprintf("nodeInfoMonitor>>GetNodeInfo %v", err), "error")
			return nil
		}
	}
	if newNodeInfo.Port == 0 {
		return errors.New("server port must > 0")
	}
	//new Node Info
	c.nodeInfo = newNodeInfo

	// Update User
	err = c.syncUserList()
	if err != nil {
		c.log(fmt.Sprintf("nodeInfoMonitor>>syncUserList %v", err), "error")
		return nil
	}

	//TODO
	//c.apiClient.GetNodeRule();

	// If nodeInfo changed
	if nodeInfoChanged {
		//disable change inbound
		nodeInfoChanged = false
	}

	return nil
}

func (c *Controller) userInfoMonitor() (err error) {
	// delay to start
	if time.Since(c.startAt) < time.Duration(c.config.UpdatePeriodic)*time.Second {
		return nil
	}

	// Get server status
	CPU, Mem, Disk, Uptime, err := serverstatus.GetSystemInfo()
	if err != nil {
		c.log(fmt.Sprintf("userInfoMonitor>>ReportNodeStatus %v", err), "error")
	} else {
		err = c.apiClient.ReportNodeStatus(
			&api.NodeStatus{
				CPU:    CPU,
				Mem:    Mem,
				Disk:   Disk,
				Uptime: Uptime,
			})
		if err != nil {
			c.log(err, "error")
		}
		c.log(fmt.Sprintf("userInfoMonitor>>ReportNodeStatus CPU %.2f, Mem %.2f, Disk %.2f, Uptime %v, err %v", CPU/100, Mem/100, Disk/100, strconv.FormatUint(Uptime, 10), err), "info")
	}

	// Get User traffic
	var userTraffic []api.UserTraffic
	var onlineUsers []api.OnlineUser
	userArr := c.author.ListUsers()
	if len(userArr) == 0 {
		c.log("traffic zero online user", "info")
		return nil
	}
	for _, user := range userArr {
		// 获取用户流量
		sent, recv := user.ResetTraffic()
		if sent == 0 && recv == 0 {
			continue
		}

		var up, down int64
		if c.nodeInfo.TrafficRate == 1 {
			up, down = sent, recv
		} else {
			up = int64(c.nodeInfo.TrafficRate * float64(sent))
			down = int64(c.nodeInfo.TrafficRate * float64(recv))
		}

		// c.log(fmt.Sprintf("get traffic user:%s sent:%d up:%d recv:%d down:%d rate:%f", user.hash, sent, up, recv, down, c.nodeInfo.TrafficRate), "debug")

		userTraffic = append(userTraffic, api.UserTraffic{
			UID:      user.UID,
			Email:    user.hash,
			Upload:   up,
			Download: down})

		// 在线用户
		iplist := user.ResetIPTable()
		for _, ip := range iplist {
			onlineUsers = append(onlineUsers, api.OnlineUser{
				UID: user.UID,
				IP:  ip,
			})

		} //END online user
	}

	userCounter := len(userTraffic)
	ipCounter := len(onlineUsers)
	//report traffic
	if userCounter > 0 {
		if err = c.apiClient.ReportUserTraffic(&userTraffic); err != nil {
			c.log(fmt.Sprintf("report traffic err:%v", err.Error()), "error")
		}
	}
	//report online users
	if ipCounter > 0 {
		if err = c.apiClient.ReportNodeOnlineUsers(&onlineUsers); err != nil {
			c.log(fmt.Sprintf("report online users err:%v", err.Error()), "error")
		}
	}
	c.log(fmt.Sprintf("userInfoMonitor>>ReportUserOnline userCounter:%d ipCounter:%d", userCounter, ipCounter), "info")

	// TODO
	// if err = c.apiClient.ReportIllegal(detectResult); err != nil {

	return nil
}

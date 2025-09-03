package shortcuts

import (
	"context"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
	"sync/atomic"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing/common/bufio"
	N "github.com/sagernet/sing/common/network"

	"github.com/sagernet/sing-box/poet/panel"
	SS "github.com/sagernet/sing-box/poet/shortcuts"

	"github.com/spf13/viper"
)

// type singPoet struct {
// }

// var instance *singPoet

// func init() {
// 	instance = new(singPoet)
// }

// func singleton() *singPoet {
// 	return instance
// }

// 获取配置文件
func getConfig(cfgFile string) *viper.Viper {
	config := viper.New()

	// Set custom path and name
	if cfgFile != "" {
		configName := path.Base(cfgFile)
		configFileExt := path.Ext(cfgFile)
		configNameOnly := strings.TrimSuffix(configName, configFileExt)
		configPath := path.Dir(cfgFile)
		config.SetConfigName(configNameOnly)
		config.SetConfigType(strings.TrimPrefix(configFileExt, "."))
		config.AddConfigPath(configPath)
	} else {
		// Set default config path
		config.SetConfigName("panel")
		config.SetConfigType("json")
		config.AddConfigPath(".")

	}

	if err := config.ReadInConfig(); err != nil {
		fmt.Println("\tErr: load json file fail> ", err)
	}

	// config.WatchConfig() // Watch the config

	return config
}

// 服务开始
func Start() error {
	ShowVersion()

	cc, _ := SS.GetObject("poetConfigPath")
	if cc == nil {
		fmt.Println("\tErr: empty panel config path !")
		os.Exit(1)
	}

	config := getConfig(cc.(string))
	panelConfig := &panel.Config{}
	if err := config.Unmarshal(panelConfig); err != nil {
		return fmt.Errorf("parse config file %v failed: %s ", cc, err)
	}

	p := panel.NewPanel(panelConfig)

	p.Start()
	defer p.Close()

	return nil

}

// 设置日志
func SetLogger(logFactory log.Factory) {
	inst := SS.Singleton()
	inst.Logger = logFactory.NewLogger("sing-poet")
}

// 设置 inbound
func SetInboud(in *adapter.Inbound, tag string) {
	SS.Singleton().SetInboud(in, tag)
}

// 设置 outbound
// func SetOutbound() {}

// 计费
func RoutedConnection(ctx context.Context, conn net.Conn, metadata adapter.InboundContext) net.Conn {
	var readCounter []*atomic.Int64
	var writeCounter []*atomic.Int64

	ss := SS.Singleton()
	contrl := ss.GetContrlWithInTag(metadata.Inbound)
	if contrl == nil {
		ss.Logger.Warn(fmt.Sprintf("RecordUserTraffic: inbound:%s contrl not found", metadata.Inbound))
		return conn
	}

	sendPtr, recvPtr, err := contrl.LoadOrCreateUserCounter(ctx, metadata)
	if err != nil {
		ss.Logger.Warn(fmt.Sprintf("RecordUserTraffic: inbound:%s user:%s counter pointer not found", metadata.Inbound, metadata.User))
		return conn
	}
	readCounter = append(readCounter, sendPtr)
	writeCounter = append(writeCounter, recvPtr)

	return bufio.NewInt64CounterConn(conn, readCounter, writeCounter)
}
func RoutedPacketConnection(ctx context.Context, conn N.PacketConn, metadata adapter.InboundContext) N.PacketConn {
	var readCounter []*atomic.Int64       //统计读取的总字节数
	var readPacketCounter []*atomic.Int64 //统计读取的数据包数量
	var writeCounter []*atomic.Int64
	var writePacketCounter []*atomic.Int64

	ss := SS.Singleton()
	contrl := ss.GetContrlWithInTag(metadata.Inbound)
	if contrl == nil {
		ss.Logger.Warn(fmt.Sprintf("RecordPacketTraffic: inbound:%s contrl not found", metadata.Inbound))
		return conn
	}

	sendPtr, recvPtr, err := contrl.LoadOrCreateUserCounter(ctx, metadata)
	if err != nil {
		ss.Logger.Warn(fmt.Sprintf("RecordPacketTraffic: inbound:%s user:%s counter pointer not found", metadata.Inbound, metadata.User))
		return conn
	}

	readCounter = append(readCounter, sendPtr)
	writeCounter = append(writeCounter, recvPtr)

	// 修改部分：只有当计数器大于0时，才创建一个值为1的新计数器并添加到包计数器切片中
	if sendPtr != nil && sendPtr.Load() > 0 {
		readCounter := &atomic.Int64{}
		readCounter.Store(1) // 设置初始值为1
		readPacketCounter = append(readPacketCounter, readCounter)
	}

	if recvPtr != nil && recvPtr.Load() > 0 {
		writeCounter := &atomic.Int64{}
		writeCounter.Store(1) // 设置初始值为1
		writePacketCounter = append(writePacketCounter, writeCounter)
	}

	return bufio.NewInt64CounterPacketConn(conn, readCounter, readPacketCounter, writeCounter, writePacketCounter)
}

// debug
func Debug(msg string) {
	SS.Singleton().Logger.Debug(msg)
}

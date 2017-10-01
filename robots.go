/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-04-01 21:11:51
 * Filename      : robot.go
 * Description   : 机器人
 * *******************************************************/
package main

import (
	"flag"
	"os"
	"os/signal"
	"protoactor/data"
	"protoactor/robots"
	"runtime"
	"syscall"
	"time"

	"github.com/golang/glog"
)

const (
	ROBOTS_NAME = "RobotMsg"
	VERSION     = "0.0.1"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//加载配置
	var config string
	flag.StringVar(&config, "conf", "./conf.json", "config path")
	flag.Parse()
	data.LoadConf(config)
	glog.Infoln("Config: ", data.Conf)
	defer glog.Flush()
	//启动远程连接服务
	rbServer := new(robots.RobotServer)
	rbServer.Addr = "127.0.0.1:7070"
	rbServer.ServerHost = "127.0.0.1"
	rbServer.ServerPort = "7011"
	rbServer.Phone = "10005000000"
	if rbServer != nil {
		rbServer.Start() //启动服务
	}
	//启动服务
	//
	closeSig := make(chan bool, 1)
	go signalProcess(closeSig) //监听关闭信号
	<-closeSig                 //通道阻塞
	//关闭服务
	//关闭websocket连接
	if rbServer != nil {
		rbServer.Close() //关闭服务
	}
	//关闭退出消息
	//
	//延迟等待
	<-time.After(10 * time.Second) //延迟关闭
}

//监听服务
func signalProcess(closeSig chan bool) {
	ch := make(chan os.Signal, 1)
	//signal.Notify(ch, syscall.SIGUSR1, syscall.SIGUSR2)
	//signal.Notify(ch, syscall.SIGHUP)
	signal.Notify(ch, os.Interrupt, os.Kill) //监听SIGINT和SIGKILL信号
	for {
		msg := <-ch
		switch msg {
		default:
			glog.Infof("close signal : %v ", msg)
			//关闭服务
			closeSig <- true
		case syscall.SIGHUP:
			//TODO
		}
	}
}

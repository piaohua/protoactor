package main

import (
	"flag"
	"os"
	"os/signal"
	"protoactor/data"
	"protoactor/process"
	"protoactor/socket"
	"runtime"
	"syscall"
	"time"

	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/golang/glog"
)

const (
	ROOMS_NAME  = "rooms"
	ONLINE_NAME = "online"
	LOGGER_NAME = "logger"
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
	//启动websocket服务
	wsServer := new(socket.WSServer)
	wsServer.Addr = "0.0.0.0:7011"
	if wsServer != nil {
		wsServer.Start() //启动服务
	}
	//启动服务
	process.RunImages()            //启动头像服务
	process.InitSupervisor()       //启动监控服务
	process.RunLogger(LOGGER_NAME) //启动日志服务
	process.RunOnline(ONLINE_NAME) //启动在线列表服务
	process.RunRooms(ROOMS_NAME)   //启动房间列表服务
	remote.Start("127.0.0.1:7072") //远程节点
	closeSig := make(chan bool, 1)
	go signalProcess(closeSig) //监听关闭信号
	<-closeSig                 //通道阻塞
	//关闭服务
	//关闭websocket连接
	if wsServer != nil {
		wsServer.Close() //关闭服务
	}
	//关闭退出消息
	process.RoomsPID.Stop()  //停止服务
	process.OnlinePID.Stop() //停止服务
	process.LogPID.Stop()    //停止服务
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

/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-04-01 21:14:34
 * Filename      : robots.go
 * Description   : 机器人
 * *******************************************************/
package robots

import (
	"log"
	"net/url"
	"protoactor/messages"
	"sync"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type RobotServer struct {
	Addr            string                  //地址
	PendingWriteNum int                     //等待写入消息长度
	MaxMsgLen       uint32                  //最大消息长度
	ServerHost      string                  //服务器地址
	ServerPort      string                  //服务器端口
	Phone           string                  //注册登录账号
	channel         chan *messages.RobotMsg //消息通道
	conns           WebsocketConnSet        //连接集合
	mutexConns      sync.Mutex              //互斥锁
	wg              sync.WaitGroup          //同步机制
}

func (server *RobotServer) Start() {
	if server.PendingWriteNum <= 0 {
		server.PendingWriteNum = 100
		glog.Infof("invalid PendingWriteNum, reset to %v", server.PendingWriteNum)
	}
	if server.MaxMsgLen <= 0 {
		server.MaxMsgLen = 10240
		glog.Infof("invalid MaxMsgLen, reset to %v", server.MaxMsgLen)
	}

	server.conns = make(WebsocketConnSet)

	if server.Addr == "" {
		panic("server.Addr empty")
	}
	// Start the remote server
	remote.Start(server.Addr)
	server.remoteRecv() //接收远程消息
}

//关闭连接
func (server *RobotServer) Close() {
	close(server.channel)

	server.mutexConns.Lock()
	for conn := range server.conns {
		conn.Close()
	}
	server.conns = nil
	server.mutexConns.Unlock()

	server.wg.Wait()
}

//接收远程消息
func (server *RobotServer) remoteRecv() {
	//create the channel
	server.channel = make(chan *messages.RobotMsg, 100) //protos中定义

	//create an actor receiving messages and pushing them onto the channel
	props := actor.FromFunc(func(context actor.Context) {
		if msg, ok := context.Message().(*messages.RobotMsg); ok {
			server.channel <- msg
		}
	})
	actor.SpawnNamed(props, "RobotMsg")

	//consume the channel just like you use to
	go func() {
		for msg := range server.channel {
			//分配机器人
			for msg.Num > 0 {
				go func(code, phone string) {
					server.runRobot(code, phone)
				}(msg.Code, server.Phone)
				server.Phone = utils.StringAdd(server.Phone)
				msg.Num--
			}
			log.Println("node msg -> ", msg)
		}
	}()
}

//启动一个机器人
func (server *RobotServer) runRobot(code, phone string) {
	host := server.ServerHost + ":" + server.ServerPort
	u := url.URL{Scheme: "ws", Host: host, Path: "/"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		glog.Errorf("robot run dial -> %v", err)
		return
	}

	server.wg.Add(1)
	defer server.wg.Done()

	server.mutexConns.Lock()
	if server.conns == nil {
		server.mutexConns.Unlock()
		conn.Close()
		return
	}
	server.conns[conn] = struct{}{}
	server.mutexConns.Unlock()

	//new robot
	robot := newRobot(conn, server.PendingWriteNum, server.MaxMsgLen)
	robot.code = code //设置邀请码
	robot.data.Phone = phone
	glog.Infof("run robot -> %s", phone)
	go robot.writePump()
	go robot.SendRegist() //发起请求,注册-登录-进入房间
	robot.readPump()

	// cleanup
	robot.Close()
	server.mutexConns.Lock()
	delete(server.conns, conn)
	server.mutexConns.Unlock()
}

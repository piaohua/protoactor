package socket

import (
	"fmt"
	"protoactor/login"
	"protoactor/messages"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
)

//test
type testWs struct{ Who string }

func (state *WSConn) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		fmt.Println("Starting, initialize actor here")
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about to shut down")
		//断开连接
		if state.playerPid != nil {
			state.playerPid.Tell(&messages.Logout{
				Sender: ctx.Self(),
			}) //通知playerPid下线
		}
	case *actor.Stopped:
		fmt.Println("Stopped, actor and its children are stopped")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about to restart")
	case *actor.ReceiveTimeout:
		glog.Infof("ReceiveTimeout: %v", ctx.Self().String())
		//TODO 超时处理(注册,登录)
		state.Close() //超时断开
	case *testWs:
		fmt.Printf("self %s\n", ctx.Self().String())
		fmt.Printf("msg.Who %v\n", msg.Who)
	case *messages.CRegist:
		//注册消息
		ctx.SetReceiveTimeout(8 * time.Second) //login timeout set
		login.Regist(msg, ctx.Self())
	case *messages.CLogin:
		//登录消息
		ctx.SetReceiveTimeout(8 * time.Second) //login timeout set
		login.Login(msg, ctx.Self())
	case *messages.RepeatLogin:
		//重复登录
		state.playerPid = nil //断开
		msg.Sender.Tell(&messages.Logined{
			Sender: msg.Repeat,
		}) //登录成功消息
		state.Close() //关闭旧连接
	case *messages.Logined:
		//登录成功
		ctx.SetReceiveTimeout(0) //login Successfully, timeout off
		state.playerPid = msg.Sender
		glog.Infof("Logined: %s", msg.Sender.String())
	case *messages.LogRegist:
		//注册日志
		msg.Ip = state.GetIPAddr()
		msg.Sender.Tell(msg)
	case *messages.LogLogin:
		//登录日志
		msg.Ip = state.GetIPAddr()
		msg.Sender.Tell(msg)
	case proto.Message:
		//消息响应
		//glog.Infof("Send: %#v", msg)
		state.Send(msg)
	default:
		glog.Infof("default: %v, %#v", ctx.Self().String(), msg)
	}
}

//初始化
func (ws *WSConn) initWs() *actor.PID {
	props := actor.FromInstance(ws) //实例
	pid := actor.Spawn(props)       //启动一个进程
	return pid
}

//pid.String()

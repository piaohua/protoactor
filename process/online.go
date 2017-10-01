package process

import (
	"fmt"
	"protoactor/messages"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/golang/glog"
)

//pid.String() == 角色ID

//全局在线列表PID
var OnlinePID *actor.PID

//test
type testOnline struct{ Who string }

//OnlineActor
type OnlineActor struct {
	OnlineSet *actor.PIDSet //在线列表, is not thread safe
}

func (state *OnlineActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		fmt.Println("Starting, initialize actor here")
		OnlinePID = ctx.Self()
		state.initOnline() //initialize
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about to shut down")
	case *actor.Stopped:
		fmt.Println("Stopped, actor and its children are stopped")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about to restart")
	case *testOnline:
		fmt.Printf("self %s\n", ctx.Self().String())
		fmt.Printf("Online %v\n", msg.Who)
	case *TestChild:
		fmt.Printf("test %s\n", ctx.Self().String())
	case *messages.OnlineReq:
		exist := state.OnlineSet.Contains(msg.Pid)
		glog.Infof("OnlineReq %s", msg.Pid.String())
		state.OnlineSet.Add(msg.Pid)
		ctx.Respond(&messages.OnlineResp{
			Result:   exist,
			Userid:   msg.Userid,
			Phone:    msg.Phone,
			Nickname: msg.Nickname,
			Sender:   msg.Sender,
		})
	}
}

//初始化
func (state *OnlineActor) initOnline() {
	state.OnlineSet = actor.NewPIDSet()
}

func newChildOnlineActor() actor.Actor {
	return &OnlineActor{}
}

//启动在线列表进程服务
func RunOnline(name string) *actor.PID {
	if name == "" {
		panic("illegal name")
	}
	pid := RunChild(name, newChildOnlineActor)
	if pid == nil {
		panic("run online failed")
	}
	return pid //监控进程
}

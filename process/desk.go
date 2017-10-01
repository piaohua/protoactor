package process

import (
	"fmt"
	"protoactor/desk"
	"protoactor/entity"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
)

//test
type testDesk struct{ Who string }

//DeskActor
type DeskActor struct {
	*entity.Desk //桌子数据
}

func (state *DeskActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		fmt.Println("Starting, initialize actor here")
		state.Pid = ctx.Self()
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about to shut down")
	case *actor.Stopped:
		fmt.Println("Stopped, actor and its children are stopped")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about to restart")
	case *actor.ReceiveTimeout:
		//超时处理(关闭)
		ctx.Self().Stop()
	case *testDesk:
		fmt.Printf("self %s\n", ctx.Self().String())
		fmt.Printf("desk %v\n", msg.Who)
	case proto.Message:
		//消息请求
		//glog.Infof("PlayerMsg: %s, %#v", ctx.Self().String(), msg)
		state.Handler(msg, ctx)
	}
}

//初始化
func InitDesk(data *entity.DeskData) *actor.PID {
	//glog.Infof("InitDesk data: %#v", data)
	props := actor.FromInstance(&DeskActor{Desk: desk.NewDesk(data)})  //桌子实例
	pid, err := actor.SpawnNamed(props, "desk"+utils.String(data.Rid)) //启动一个进程
	if err != nil {
		fmt.Printf("init desk err -> %v", err)
		return pid //TODO name exists
	}
	return pid
}

//pid.String() == "desk"+牌桌ID

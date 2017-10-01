package process

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/maps/hashmap"
)

//pid.String() == "room"+邀请码

//全局房间列表PID
var RoomsPID *actor.PID

//test
type testRooms struct{ Who string }

//RoomsActor
type RoomsActor struct {
	RoomLastID uint32        //房间唯一ID,递增
	RoomsSet   *actor.PIDSet //已经开启房间进程PID, is not thread safe
	RoomsMap   *hashmap.Map  //房间列表, is not thread safe
}

func (state *RoomsActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		fmt.Println("Starting, initialize actor here")
		RoomsPID = ctx.Self()
		state.initRooms() //initialize
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about to shut down")
	case *actor.Stopped:
		fmt.Println("Stopped, actor and its children are stopped")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about to restart")
	case *testRooms:
		fmt.Printf("self %s\n", ctx.Self().String())
		fmt.Printf("Rooms %v\n", msg.Who)
	case *TestChild:
		fmt.Printf("test %s\n", ctx.Self().String())
	default:
		state.Handler(msg, ctx)
	}
}

//初始化
func (state *RoomsActor) initRooms() {
	state.RoomLastID = 1 //TODO
	state.RoomsSet = actor.NewPIDSet()
	state.RoomsMap = hashmap.New()
}

func newChildRoomsActor() actor.Actor {
	return &RoomsActor{}
}

//启动房间列表进程服务
func RunRooms(name string) *actor.PID {
	if name == "" {
		panic("illegal name")
	}
	pid := RunChild(name, newChildRoomsActor)
	if pid == nil {
		panic("run rooms failed")
	}
	return pid //监控进程
}

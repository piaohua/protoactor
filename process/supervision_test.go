package process

import (
	"fmt"
	"testing"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
)

func TestSupervisor(t *testing.T) {
	InitSupervisor()
	//创建一个一对一监控服务
	props := actor.
		FromProducer(newParentActor).
		WithSupervisor(Supervisor) //监控进程定义
	pid := actor.Spawn(props)              //新启一个监控进程
	pid.Tell(&testSuper{p: newChildActor}) //向监控进程发一个消息
	t.Log(pid.String())
	pid2 := actor.Spawn(props)              //新启一个监控进程
	pid2.Tell(&testSuper{p: newChildActor}) //向监控进程发一个消息
	t.Log(pid2.String())
	<-time.After(5 * time.Second)
	//
	console.ReadLine()
}

type childActor struct{}

func (state *childActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Starting, initialize actor here")
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about to shut down")
		<-time.After(3 * time.Second)
	case *actor.Stopped:
		fmt.Println("Stopped, actor and its children are stopped")
		<-time.After(2 * time.Second)
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about to restart")
	case *TestChild:
		fmt.Printf("Hello %v\n", msg.Who)
		//panic("Ouch")
		context.Self().Stop()
	}
}

func newChildActor() actor.Actor {
	return &childActor{}
}

package process

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//test
type testSuper struct{ p actor.Producer }

//test
type TestChild struct{ Who string }

//启动子进程
type StartChild struct {
	Name string         //字进程名字
	P    actor.Producer //定义
}

//ParentActor
type ParentActor struct{}

func (state *ParentActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *testSuper: //from super process
		props := actor.FromProducer(msg.p)       //子进程定义
		child := ctx.Spawn(props)                //启动一个子进程
		child.Tell(&TestChild{Who: "testChild"}) //向子进程发送一个消息
	case *StartChild:
		props := actor.FromProducer(msg.P)            //子进程定义
		child, err := ctx.SpawnNamed(props, msg.Name) //启动一个子进程
		if err != nil {
			panic(fmt.Sprintf("start child failed Name:%s, err:%s", msg.Name, err))
		}
		child.Tell(&TestChild{Who: "testChild"}) //向子进程发送一个消息
	}
}

//新启一个ParentActor
func newParentActor() actor.Actor {
	return &ParentActor{}
}

var Supervisor actor.SupervisorStrategy

//初始化
func InitSupervisor() {
	//一对一监控,1000时间内最大重试次数10次
	Supervisor = actor.NewOneForOneStrategy(10, 1000, decider)
	if Supervisor == nil {
		panic("Supervisor Equal NIL")
	}
	//一对多监控,1000时间内最大重试次数10次
	//Supervisors = actor.NewAllForOneStrategy(10, 1000, decider)
}

//子进程失败处理函数
func decider(reason interface{}) actor.Directive {
	fmt.Println("handling failure for child")
	return actor.StopDirective
}

//启动一个子进程
func RunChild(name string, p actor.Producer) *actor.PID {
	if name == "" {
		panic("illegal name")
	}
	//创建一个一对一监控服务
	props := actor.
		FromProducer(newParentActor).
		WithSupervisor(Supervisor) //监控进程定义
	pid, err := actor.SpawnNamed(props, name) //启动一个监控进程
	if err != nil {
		fmt.Printf("run child err -> %v", err)
		return nil
	}
	pid.Tell(&StartChild{
		Name: name, //房间列表进程名字
		P:    p,
	})
	return pid //监控进程
}

package process

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//pid.String() == "logger"

//全局日志PID
var LogPID *actor.PID

//test
type testLogger struct{ Who string }

//LoggerActor
type LoggerActor struct {
	//mailbox
}

func (state *LoggerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		fmt.Println("Starting, initialize actor here")
		LogPID = ctx.Self()
		state.initLogger() //initialize
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about to shut down")
	case *actor.Stopped:
		fmt.Println("Stopped, actor and its children are stopped")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about to restart")
	case *testLogger:
		fmt.Printf("self %s\n", ctx.Self().String())
		fmt.Printf("Logger %v\n", msg.Who)
	case *TestChild:
		fmt.Printf("test %s\n", ctx.Self().String())
	default:
		state.Handler(msg, ctx)
	}
}

//初始化
func (state *LoggerActor) initLogger() {
}

func newChildLoggerActor() actor.Actor {
	return &LoggerActor{}
}

//启动日志进程服务
func RunLogger(name string) *actor.PID {
	if name == "" {
		panic("illegal name")
	}
	pid := RunChild(name, newChildLoggerActor)
	if pid == nil {
		panic("run logger failed")
	}
	return pid //监控进程
}

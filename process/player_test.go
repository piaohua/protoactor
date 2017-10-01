package process

import (
	"fmt"
	"testing"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
)

func TestPlayer(t *testing.T) {
	pid1 := InitPlayer("1")
	if pid1 == nil {
		panic("pid equal nil")
	}
	pid1.Tell(&testPlayer{Who: "Roger"})
	//pid1.Stop()
	<-time.After(time.Duration(2) * time.Second)
	//pid1.Tell(&testPlayer{Who: "Roger"})
	//pid name exists
	pid2 := InitPlayer("1")
	pid2.Tell(&testPlayer{Who: "Jackie"})
	//new pid
	pid3 := InitPlayer("2")
	pid3.Tell(&testPlayer{Who: "Amy"})
	console.ReadLine()
}

func TestLoginPlayer(t *testing.T) {
	props1 := actor.FromInstance(&testLoginActor{})
	pid1 := actor.Spawn(props1)
	<-time.After(time.Duration(2) * time.Second)
	pid1.Tell(&userLogin{Who: "Roger1"})
	//--
	props2 := actor.FromInstance(&testLoginActor{})
	pid2 := actor.Spawn(props2)
	time.Sleep(6 * time.Second)
	//<-time.After(time.Duration(5) * time.Second)
	pid2.Tell(&userLogin{Who: "Roger2"})
	//--
	props3 := actor.FromInstance(&testLoginActor{})
	pid3 := actor.Spawn(props3)
	<-time.After(time.Duration(6) * time.Second)
	pid3.Tell(&userLogin{Who: "Roger3"})
	//--
	props4 := actor.FromInstance(&testLoginActor{})
	pid4 := actor.Spawn(props4)
	resp := make(chan bool, 1)
	pid4.Tell(&userLoginResp{resp: resp})
	t.Log("resp chan ", <-resp)
	console.ReadLine()
}

type userLogin struct{ Who string }
type userLoginResp struct{ resp chan bool }
type testLoginActor struct {
	Who   string
	num   int
	login bool
}

func (state *testLoginActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		fmt.Printf("Started %v\n", ctx.Self().String())
		ctx.SetReceiveTimeout(5 * time.Second) //timeout set
	case *actor.ReceiveTimeout:
		state.num++
		fmt.Printf("ReceiveTimeout %s, num %d\n", ctx.Self().String(), state.num)
		ctx.Self().Stop() //timeout stoped
	case *userLogin:
		fmt.Printf("self %s\n", ctx.Self().String())
		fmt.Printf("test %v\n", msg.Who)
		state.Who = msg.Who
		state.login = true
		fmt.Printf("state %v\n", state)
		ctx.SetReceiveTimeout(0) //login Successfully, timeout off
	case *userLoginResp:
		fmt.Printf("self %s\n", ctx.Self().String())
		time.Sleep(3 * time.Second)
		close(msg.resp)
	}
}

package process

import (
	"fmt"
	"testing"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
)

func TestOnline(t *testing.T) {
	var OnlineSet *actor.PIDSet
	OnlineSet = actor.NewPIDSet()
	<-time.After(2 * time.Second)
	userid := "1000"
	pid := actor.NewLocalPID(userid)
	OnlineSet.Add(pid)
	//--
	props2 := actor.FromInstance(&testOnlineActor{})
	pid2 := actor.Spawn(props2)
	pid2.Tell(&userOnline{Who: "Roger1"})
	pid2.Tell(&userOnline{Who: "Roger2"})
	pid2.Tell(&userOnline{Who: "Roger3"})
	OnlineSet.Add(pid2)
	//--
	props3 := actor.FromInstance(&testOnlineActor{})
	pid3, _ := actor.SpawnNamed(props3, "1001")
	pid3.Tell(&userOnline{Who: "Diane1"})
	pid3.Tell(&userOnline{Who: "Diane2"})
	pid3.Tell(&userOnline{Who: "Diane3"})
	pid4, err := actor.SpawnNamed(props3, "1001")
	t.Log(err)
	pid5 := actor.NewLocalPID("1001")
	pid6 := actor.NewLocalPID("1002")
	t.Log(pid3.String(), pid4.String(), pid5.String(), pid6.String(), pid3 == pid4)
	pid5.Tell(&userOnline{Who: "Annie"})
	pid6.Tell(&userOnline{Who: "Jackie"})
	OnlineSet.Add(pid3)
	<-time.After(2 * time.Second)
	pid3.Stop()
	<-time.After(2 * time.Second)
	t.Log("pid3 -> ", OnlineSet.Contains(pid3), pid3.Size())
	pid3.Tell(&userOnline{Who: "Diane4"})
	//---
	t.Log(OnlineSet.Contains(pid))
	t.Log(OnlineSet.Len())
	t.Log(OnlineSet.Empty())
	t.Log(OnlineSet.Values())
	t.Log(OnlineSet.Remove(pid))
	t.Log(OnlineSet.Remove(pid))
	<-time.After(2 * time.Second)
	//
	console.ReadLine()
}

type userOnline struct{ Who string }
type testOnlineActor struct{ Who string }

func (state *testOnlineActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		fmt.Printf("Started %v\n", ctx.Self().String())
		ctx.SetReceiveTimeout(1 * time.Second)
	case *actor.ReceiveTimeout:
		fmt.Printf("ReceiveTimeout %s\n", ctx.Self().String())
	case *userOnline:
		fmt.Printf("self %s\n", ctx.Self().String())
		fmt.Printf("test %v\n", msg.Who)
		state.Who = msg.Who
		fmt.Printf("state %v\n", state)
	}
}

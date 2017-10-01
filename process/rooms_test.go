package process

import (
	"fmt"
	"testing"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/emirpasic/gods/sets/hashset"
)

func TestRooms(t *testing.T) {
	var RoomsSet *actor.PIDSet
	RoomsSet = actor.NewPIDSet()
	<-time.After(2 * time.Second)
	code := "123456"
	pid := actor.NewLocalPID(code)
	RoomsSet.Add(pid)
	//--
	props2 := actor.FromInstance(&testRoomsActor{})
	pid2 := actor.Spawn(props2)
	pid2.Tell(&userRooms{Who: "Roger"})
	RoomsSet.Add(pid2)
	props3 := actor.FromInstance(&testRoomsActor{})
	pid3, _ := actor.SpawnNamed(props3, "123457")
	pid3.Tell(&userRooms{Who: "Diane"})
	RoomsSet.Add(pid3)
	<-time.After(2 * time.Second)
	//---
	t.Log(RoomsSet.Contains(pid))
	t.Log(RoomsSet.Len())
	t.Log(RoomsSet.Empty())
	t.Log(RoomsSet.Values())
	t.Log(RoomsSet.Remove(pid))
	t.Log(RoomsSet.Remove(pid))
	<-time.After(2 * time.Second)
	//
	console.ReadLine()
}

type userRooms struct{ Who string }
type testRoomsActor struct{ Who string }

func (state *testRoomsActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *userRooms:
		fmt.Printf("self %s\n", ctx.Self().String())
		fmt.Printf("test %v\n", msg.Who)
		state.Who = msg.Who
		fmt.Printf("state %v\n", state)
	}
}

func TestSetRooms(t *testing.T) {
	rooms := hashset.New()
	<-time.After(2 * time.Second)
	code := "123456"
	pid1 := actor.NewLocalPID(code)
	rooms.Add(pid1)
	//--
	props2 := actor.FromInstance(&testRoomActor{})
	pid2 := actor.Spawn(props2)
	pid2.Tell(&userRoom{Who: "Roger"})
	rooms.Add(pid2)
	//--
	props3 := actor.FromInstance(&testRoomActor{})
	pid3, _ := actor.SpawnNamed(props3, "123457")
	pid3.Tell(&userRoom{Who: "Diane"})
	rooms.Add(pid3)
	//--
	for _, tmp := range rooms.Values() {
		room := tmp.(*actor.PID)
		room.Tell(&userRoom{Who: "Annie"})
	}
	<-time.After(2 * time.Second)
	//---
	t.Log(rooms.Contains(pid1))
	t.Log(rooms.Size())
	t.Log(rooms.Empty())
	t.Log(rooms.Values())
	rooms.Remove(pid1)
	rooms.Remove(pid1)
	t.Log(rooms.Values())
	<-time.After(2 * time.Second)
	//
	console.ReadLine()
}

type userRoom struct{ Who string }
type testRoomActor struct{ Who string }

func (state *testRoomActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *userRoom:
		fmt.Printf("self %s\n", ctx.Self().String())
		fmt.Printf("test %v\n", msg.Who)
		state.Who = msg.Who
		fmt.Printf("state %v\n", state)
	}
}

package main

import (
	"fmt"

	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
)

type hello struct{ Who string }
type helloActor struct{}

func (state *helloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *hello:
		fmt.Printf("Hello %v\n", msg.Who)
	}
}

func main() {
	props := actor.FromInstance(&helloActor{})
	pid := actor.Spawn(props)
	pid.Tell(&hello{Who: "Roger"})
	//
	props2 := actor.FromInstance(&helloActor2{})
	pid2 := actor.Spawn(props2)
	pid2.Request("test", pid)
	console.ReadLine()
}

type hello2 struct{ Who string }
type helloActor2 struct{}

func (state *helloActor2) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *hello2:
		fmt.Printf("Hello %v\n", msg.Who)
	case string:
		fmt.Printf("string %v\n", msg)
		context.Respond(&hello{Who: "Jackie"})
		context.Respond(&hello{Who: "Jackie2"})
	}
}

package process

import (
	"fmt"
	"testing"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func TestNode1(t *testing.T) {
	node1()
	console.ReadLine()
}

func node1() {
	timeout := 5 * time.Second
	remote.Start("127.0.0.1:8081")

	props := actor.FromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		case *actor.Started:
			fmt.Println("Local actor started")
			pid, err := remote.SpawnNamed("127.0.0.1:8080", "myRemote", "remote", timeout)
			if err != nil {
				fmt.Println("Local failed to spawn remote actor")
				return
			}
			fmt.Println("Local spawned remote actor")
			ctx.Watch(pid)
			fmt.Println("Local is watching remote actor")
		case *actor.Terminated:
			fmt.Printf("Local got terminated message %+v", msg)
		}
	})
	actor.Spawn(props)
	console.ReadLine()
}

func TestNode2(t *testing.T) {
	node2()
	console.ReadLine()
}

func node2() {
	//empty actor just to have something to remote spawn
	props := actor.FromFunc(func(ctx actor.Context) {})
	remote.Register("remote", props)

	remote.Start("127.0.0.1:8080")

	console.ReadLine()
}

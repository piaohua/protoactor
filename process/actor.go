package process

import (
	"fmt"
	"log"
	"strconv"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/router"
)

type hello struct{ Who string }
type helloActor struct{}

func (state *helloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *hello:
		fmt.Printf("self %s\n", context.Self().String())
		fmt.Printf("Hello %v\n", msg.Who)
	}
}

type hello2 struct {
	Who string
	pid *actor.PID
}
type helloActor2 struct{}

func (state *helloActor2) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *hello2:
		fmt.Printf("self 2 %s\n", context.Self().String())
		fmt.Printf("Hello 2 %v, %s\n", msg.Who, msg.pid.String())
		msg.pid.Request(&hello{Who: "Diane"}, context.Self())
	}
}

func actor_run() {
	props := actor.FromInstance(&helloActor{})
	pid := actor.Spawn(props)
	pid.Tell(&hello{Who: "Roger"})
	//pid.Stop()
	<-time.After(time.Duration(2) * time.Second)
	pid.Tell(&hello{Who: "Roger"})
	//-----
	props2 := actor.FromInstance(&helloActor2{})
	pid2 := actor.Spawn(props2)
	pid2.Tell(&hello2{Who: "Jackie", pid: pid})
	<-time.After(time.Duration(5) * time.Second)
	//pid2.Stop()
	console.ReadLine()
	<-time.After(time.Duration(2) * time.Second)

	//routing()
}

type myMessage struct{ i int }

func (m *myMessage) Hash() string {
	return strconv.Itoa(m.i)
}

func routing() {

	log.Println("Round robin routing:")
	act := func(context actor.Context) {
		switch msg := context.Message().(type) {
		case *myMessage:
			log.Printf("%v got message %d", context.Self(), msg.i)
		}
	}

	pid := actor.Spawn(router.NewRoundRobinPool(5).WithFunc(act))
	for i := 0; i < 10; i++ {
		pid.Tell(&myMessage{i})
	}
	time.Sleep(1 * time.Second)
	log.Println("Random routing:")
	pid = actor.Spawn(router.NewRandomPool(5).WithFunc(act))
	for i := 0; i < 10; i++ {
		pid.Tell(&myMessage{i})
	}
	time.Sleep(1 * time.Second)
	log.Println("ConsistentHash routing:")
	pid = actor.Spawn(router.NewConsistentHashPool(5).WithFunc(act))
	for i := 0; i < 10; i++ {
		pid.Tell(&myMessage{i})
	}
	time.Sleep(1 * time.Second)
	log.Println("BroadcastPool routing:")
	pid = actor.Spawn(router.NewBroadcastPool(5).WithFunc(act))
	for i := 0; i < 10; i++ {
		pid.Tell(&myMessage{i})
	}
	console.ReadLine()
}

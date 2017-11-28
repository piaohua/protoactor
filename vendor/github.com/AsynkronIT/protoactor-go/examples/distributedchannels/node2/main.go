package main

import (
	"fmt"
	"log"
	"time"
	"runtime"

	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/examples/distributedchannels/messages"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	remote.Start("127.0.0.1:8080")
	//create the channel
	channel := make(chan *messages.MyMessage)

	//create an actor receiving messages and pushing them onto the channel
	props := actor.FromFunc(func(context actor.Context) {
		if msg, ok := context.Message().(*messages.MyMessage); ok {
			channel <- msg
		}
	})
	actor.SpawnNamed(props, "MyMessage")

	//consume the channel just like you use to
	go func() {
		for msg := range channel {
			log.Println("node2 msg -> ", msg)
		}
	}()

	go remote_recv2()

	go remote_send2()

	console.ReadLine()
}

func remote_recv2() {
	//create the channel2
	channel2 := make(chan *messages.MyMessage)

	//create an actor receiving messages and pushing them onto the channel2
	props := actor.FromFunc(func(context actor.Context) {
		if msg, ok := context.Message().(*messages.MyMessage); ok {
			channel2 <- msg
		}
	})
	actor.SpawnNamed(props, "MyMessage2")

	//consume the channel2 just like you use to
	go func() {
		for msg := range channel2 {
			log.Println("node2 msg2 -> ", msg)
		}
	}()
}

//----------
func newMyMessageSenderChannel2() (chan<- *messages.MyMessage, *actor.PID) {
	channel2 := make(chan *messages.MyMessage)
	remote := actor.NewPID("127.0.0.1:8088", "MyMessage2")
	go func() {
		for msg := range channel2 {
			remote.Tell(msg)
		}
	}()

	return channel2, remote
}

func remote_send2() {
	time.Sleep(10 * time.Second) //等待node1启动
	channel2, remote := newMyMessageSenderChannel2()
	str := remote.String()

	for i := 0; i < 2; i++ {
		message := &messages.MyMessage{
			Message: fmt.Sprintf("hello, This is Jackie. %v, %v", i, str),
		}
		channel2 <- message
	}
}

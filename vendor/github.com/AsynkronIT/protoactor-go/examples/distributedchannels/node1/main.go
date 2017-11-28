package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/examples/distributedchannels/messages"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func newMyMessageSenderChannel() chan<- *messages.MyMessage {
	channel := make(chan *messages.MyMessage)
	remote := actor.NewPID("127.0.0.1:8080", "MyMessage")
	go func() {
		for msg := range channel {
			remote.Tell(msg)
		}
	}()

	return channel
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	remote.Start("127.0.0.1:8088")
	channel := newMyMessageSenderChannel()

	for i := 0; i < 3; i++ {
		message := &messages.MyMessage{
			Message: fmt.Sprintf("hello %v", i),
		}
		channel <- message
	}

	go remote_send2()

	go remote_recv2()

	console.ReadLine()
}

func newMyMessageSenderChannel2() chan<- *messages.MyMessage {
	channel2 := make(chan *messages.MyMessage)
	remote := actor.NewPID("127.0.0.1:8080", "MyMessage2")
	go func() {
		for msg := range channel2 {
			remote.Tell(msg)
		}
	}()

	return channel2
}

func remote_send2() {
	channel2 := newMyMessageSenderChannel2()

	for i := 0; i < 2; i++ {
		message := &messages.MyMessage{
			Message: fmt.Sprintf("hello, This is Diane. %v", i),
		}
		channel2 <- message
	}
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
			log.Println("node1 msg2 -> ", msg)
		}
	}()
}

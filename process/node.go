package process

import (
	"fmt"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/examples/remoterouting/messages"
	"github.com/AsynkronIT/protoactor-go/mailbox"
	"github.com/AsynkronIT/protoactor-go/remote"
)

type testNode struct{ Who string }

type RemoteActor struct {
	name  string
	count int
}

func (a *RemoteActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *messages.Ping:
		ctx.Respond(&messages.Pong{})
	}
}

func newRemoteActor(name string) actor.Producer {
	return func() actor.Actor {
		return &RemoteActor{
			name: name,
		}
	}
}

func newRemoteActor2() actor.Actor {
	return &RemoteActor{}
}

func NewRemote(bind, name string) {
	remote.Start(bind)
	props := actor.
		FromProducer(newRemoteActor(name)).
		WithMailbox(mailbox.Bounded(10000))
	actor.SpawnNamed(props, "remote")
	remote.Register("hello", actor.FromProducer(newRemoteActor2))
}

func activate() {
	timeout := 5 * time.Second
	pid, _ := remote.SpawnNamed("127.0.0.1:8080", "remote", "hello", timeout)
	res, _ := pid.RequestFuture(&messages.Ping{}, timeout).Result()
	response := res.(*messages.Pong)
	fmt.Println(response)
}

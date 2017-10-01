package process

import (
	"testing"
	"time"

	console "github.com/AsynkronIT/goconsole"
)

func TestDesk(t *testing.T) {
	pid := InitDesk("1")
	if pid == nil {
		panic("pid equal nil")
	}
	pid.Tell(&testDesk{Who: "Roger"})
	pid.Stop()
	<-time.After(time.Duration(2) * time.Second)
	pid.Tell(&testDesk{Who: "Roger"})
	//
	console.ReadLine()
}

package process

import (
	"testing"

	console "github.com/AsynkronIT/goconsole"
)

func TestImageServer(t *testing.T) {
	RunImages()
	console.ReadLine()
}

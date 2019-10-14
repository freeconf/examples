package basics_test

import (
	"testing"

	"github.com/freeconf/examples/basics"
)

func TestCar(t *testing.T) {
	c := basics.New()
	c.Speed = 1
	update := make(chan bool)
	c.OnUpdate(func(c *basics.Car) {
		update <- true
	})
	c.Start()
	<-update
	if !c.Running {
		t.Error("not running when starting")
	}
}

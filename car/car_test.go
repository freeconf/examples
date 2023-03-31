package car

import (
	"fmt"
	"testing"
	"time"

	"github.com/freeconf/yang/fc"
)

func TestCar(t *testing.T) {
	c := New()
	c.PollInterval = time.Millisecond
	events := make(chan updateEvent)
	c.OnUpdate(func(e updateEvent) {
		fmt.Printf("got event %s\n", e)
		events <- e
	})
	c.Start()
	t.Log("waiting for car events...")
	fc.AssertEqual(t, carStarted, <-events)
	fc.AssertEqual(t, flatTire, <-events)
	fc.AssertEqual(t, carStopped, <-events)
	c.replaceTires()
	fc.AssertEqual(t, carStarted, <-events)
}

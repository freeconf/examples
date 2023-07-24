package car

import (
	"fmt"
	"testing"
	"time"

	"github.com/freeconf/yang/fc"
)

// Quick test of car's features using direct access to fields and methods
func TestCar(t *testing.T) {
	c := New()
	c.PollInterval = time.Millisecond
	c.Speed = 1000

	events := make(chan updateEvent)
	unsub := c.OnUpdate(func(e updateEvent) {
		fmt.Printf("got event %s\n", e)
		events <- e
	})
	t.Log("waiting for car events...")
	c.Start()

	fc.AssertEqual(t, carStarted, <-events)
	fc.AssertEqual(t, flatTire, <-events)
	fc.AssertEqual(t, carStopped, <-events)
	c.replaceTires()
	c.Start()

	fc.AssertEqual(t, carStarted, <-events)
	unsub.Close()
	c.Stop()
}

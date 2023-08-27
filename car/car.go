package car

import (
	"container/list"
	"fmt"
	"math/rand"
	"time"
)

// ////////////////////////
// C A R - example application
//
// This has nothing to do with FreeCONF, just an example application written in Go.
// that models a running car that can get flat tires when tires are worn.
type Car struct {
	Tire []*Tire

	Miles   float64
	Running bool

	// When the tires were last rotated
	LastRotation int64

	// Default speed value is in yang model file and free's your code
	// from hardcoded values, even if they are only default values
	// units milliseconds/mile
	Speed int

	// How fast to apply speed. Default it 1s or "miles per second"
	PollInterval time.Duration

	// Listeners are common on manageable code.  Having said that, listeners
	// remain relevant to your application.  The manage.go file is responsible
	// for bridging the conversion from application to management api.
	listeners *list.List
}

// CarListener for receiving car update events
type CarListener func(UpdateEvent)

// car event types
type UpdateEvent int

const (
	CarStarted UpdateEvent = iota + 1
	CarStopped
	FlatTire
)

func (e UpdateEvent) String() string {
	strs := []string{
		"unknown",
		"carStarted",
		"carStopped",
		"flatTire",
	}
	if int(e) < len(strs) {
		return strs[e]
	}
	return "invalid"
}

func New() *Car {
	c := &Car{
		listeners:    list.New(),
		Speed:        1000,
		PollInterval: time.Second,
	}
	c.NewTires()
	return c
}

// Stop will take up to poll_interval seconds to come to a stop
func (c *Car) Stop() {
	c.Running = false
}

func (c *Car) Start() {
	if c.Running {
		return
	}
	go func() {
		c.Running = true
		c.updateListeners(CarStarted)
		defer func() {
			c.Running = false
			c.updateListeners(CarStopped)
		}()
		for c.Speed > 0 {
			poll := time.NewTicker(c.PollInterval)
			for c.Running {
				for range poll.C {
					c.Miles += float64(c.Speed)
					for _, t := range c.Tire {
						t.endureMileage(c.Speed)
						if t.Flat {
							c.updateListeners(FlatTire)
							return
						}
					}
				}
			}
		}
	}()
}

// OnUpdate to listen for car events like start, stop and flat tire
func (c *Car) OnUpdate(l CarListener) Subscription {
	return NewSubscription(c.listeners, c.listeners.PushBack(l))
}

func (c *Car) NewTires() {
	c.Tire = make([]*Tire, 4)
	c.LastRotation = int64(c.Miles)
	for pos := 0; pos < len(c.Tire); pos++ {
		c.Tire[pos] = &Tire{
			Pos:  pos,
			Wear: 100,
			Size: "H15",
		}
	}
}

func (c *Car) ReplaceTires() {
	for _, t := range c.Tire {
		t.Replace()
	}
	c.LastRotation = int64(c.Miles)
}

func (c *Car) RotateTires() {
	x := c.Tire[0]
	c.Tire[0] = c.Tire[1]
	c.Tire[1] = c.Tire[2]
	c.Tire[2] = c.Tire[3]
	c.Tire[3] = x
	for i, t := range c.Tire {
		t.Pos = i
	}
	c.LastRotation = int64(c.Miles)
}

func (c *Car) updateListeners(e UpdateEvent) {
	fmt.Printf("car %s\n", e)
	i := c.listeners.Front()
	for i != nil {
		i.Value.(CarListener)(e)
		i = i.Next()
	}
}

// T I R E
type Tire struct {
	Pos  int
	Size string
	Flat bool
	Wear float64
	Worn bool
}

func (t *Tire) Replace() {
	t.Wear = 100
	t.Flat = false
	t.Worn = false
}

func (t *Tire) checkIfFlat() {
	if !t.Flat {
		// emulate that the more wear a tire has, the more likely it will
		// get a flat, but there is always a chance.
		t.Flat = (t.Wear - (rand.Float64() * 10)) < 0
	}
}

func (t *Tire) endureMileage(speed int) {
	// Wear down [0.0 - 0.5] of each tire proportionally to the tire position
	t.Wear -= (float64(speed) / 100) * float64(t.Pos) * rand.Float64()
	t.checkIfFlat()
	t.checkForWear()
}

func (t *Tire) checkForWear() bool {
	return t.Wear < 20
}

///////////////////////
// U T I L

// Subscription is handle into a list.List that when closed
// will automatically remove item from list.  Useful for maintaining
// a set of listeners that can easily remove themselves.
type Subscription interface {
	Close() error
}

// NewSubscription is used by subscription managers to give a token
// to caller the can close to unsubscribe to events
func NewSubscription(l *list.List, e *list.Element) Subscription {
	return &listSubscription{l, e}
}

type listSubscription struct {
	l *list.List
	e *list.Element
}

// Close will unsubscribe to events.
func (sub *listSubscription) Close() error {
	sub.l.Remove(sub.e)
	return nil
}

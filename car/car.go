package car

import (
	"container/list"
	"math/rand"
	"time"
)

// ////////////////////////
// C A R
// Your application code.
//
// Notice there are no reference to FreeCONF in this file.  This means your
// code remains:
// - unit test-able
// - Not auto-generated from model files
// - free of golang source code annotations/tags.
type Car struct {
	Tire []*Tire

	// Not everything has to be structs, using a map may be useful
	// in early prototyping
	Specs map[string]interface{}

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

type CarListener func(updateEvent)

type updateEvent int

const (
	carStarted updateEvent = iota + 1
	carStopped
	flatTire
)

func (e updateEvent) String() string {
	strs := []string{
		"unknown",
		"started",
		"stopped",
		"flat",
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
		Specs:        make(map[string]interface{}),
		PollInterval: time.Second,
	}
	c.newTires()
	return c
}

func (c *Car) newTires() {
	c.Tire = make([]*Tire, 4)
	c.LastRotation = int64(c.Miles)
	for pos := 0; pos < len(c.Tire); pos++ {
		c.Tire[pos] = &Tire{
			Pos:  pos,
			Wear: 100,
		}
	}
}

func (c *Car) Start() {
	if c.Running {
		return
	}
	go func() {
		c.Running = true
		c.updateListeners(carStarted)
		for c.Speed > 0 {
			poll := time.NewTicker(c.PollInterval)
			for c.Running {
				for range poll.C {
					c.Miles += float64(c.Speed)

					for _, t := range c.Tire {

						t.endureMileage(c.Speed)

						if t.Flat {
							c.updateListeners(flatTire)
							goto done
						}
					}
				}
			}
		}
	done:
		c.Running = false
		c.updateListeners(carStopped)
	}()
}

func (c *Car) OnUpdate(l CarListener) Subscription {
	return NewSubscription(c.listeners, c.listeners.PushBack(l))
}

func (c *Car) updateListeners(e updateEvent) {
	i := c.listeners.Front()
	for i != nil {
		i.Value.(CarListener)(e)
		i = i.Next()
	}
}

func (c *Car) replaceTires() {
	for _, t := range c.Tire {
		t.replace()
	}
	c.LastRotation = int64(c.Miles)
	c.Start()
}

func (c *Car) rotateTires() {
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

// T I R E
type Tire struct {
	Pos  int
	Size string
	Flat bool
	Wear float64
	Worn bool
}

func (t *Tire) replace() {
	t.Wear = 100
	t.Flat = false
	t.Worn = false
}

func (t *Tire) checkIfFlat() {
	if !t.Flat {
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

package main

import (
	"container/list"
	"os"
	"time"

	"github.com/freeconf/gconf/meta"
	"github.com/freeconf/gconf/node"

	"github.com/freeconf/gconf/device"
	"github.com/freeconf/gconf/meta/yang"
	"github.com/freeconf/gconf/nodes"
	"github.com/freeconf/gconf/restconf"
)

func main() {

	// Your app
	car := NewCar()

	// Add management
	uiPath := &meta.FileStreamSource{Root: "."}
	d := device.NewWithUi(yang.YangPath(), uiPath)
	if err := d.Add("car", manage(car)); err != nil {
		panic(err)
	}
	restconf.NewServer(d)
	if err := d.ApplyStartupConfig(os.Stdin); err != nil {
		panic(err)
	}
	car.Refuel(18)

	// trick to sleep forever...
	select {}
}

type Car struct {
	Speed     int
	Miles     int64
	Gas       int
	State     int
	listeners *list.List
}

type listener func(int)

func NewCar() *Car {
	return &Car{
		listeners: list.New(),
	}
}

func (c *Car) update(state int) {
	c.State = state
	for p := c.listeners.Front(); p != nil; p = p.Next() {
		p.Value.(listener)(state)
	}
}

func (c *Car) onUpdate(l listener) *list.Element {
	return c.listeners.PushBack(l)
}

var Mpg = int64(30)

const (
	OutOfGas = iota
	Running
)

func (c *Car) Start() {
	if c.State == Running {
		return
	}
	c.update(Running)
	defer c.update(OutOfGas)
	for c.Gas > 0 {
		<-time.After(time.Duration(c.Speed) * time.Millisecond)
		c.Miles++
		if (c.Miles % Mpg) == 0 {
			c.Gas--
		}
	}
}

func (c *Car) Refuel(gallons int) {
	c.Gas += gallons
	if c.State != Running {
		go c.Start()
	}
}

// implement your mangement api
func manage(c *Car) node.Node {

	// this is one of many options to reuse management code
	return &nodes.Extend{

		// use reflect when possible
		Base: nodes.ReflectChild(c),

		// handle 'rpc' or 'action' in YANG
		OnAction: func(parent node.Node, req node.ActionRequest) (node.Node, error) {
			switch req.Meta.Ident() {
			case "refuel":
				v, err := req.Input.Get("gas")
				if err != nil {
					return nil, err
				}
				c.Refuel(v.(int))
			}
			return nil, nil
		},

		// handle 'notification' in YANG
		OnNotify: func(parent node.Node, req node.NotifyRequest) (node.NotifyCloser, error) {
			switch req.Meta.Ident() {
			case "update":

				// See c2.NewSubscription for simple implementation of managing listeners
				// like this
				e := c.onUpdate(func(update int) {
					req.Send(manage(c))
				})
				unsub := func() error {
					c.listeners.Remove(e)
					return nil
				}

				return unsub, nil
			}
			return nil, nil
		},
	}
}

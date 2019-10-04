package car

import (
	"fmt"

	"github.com/freeconf/yang/node"
)

// This is where you can implement any customizations and required function
// implementations

// Start function was written by hand
func (c *Car) doStart(r node.ActionRequest) (node.Node, error) {
	c.Running = true
	return nil, nil
}

func (c *Car) onUpdate(r node.NotifyRequest) (node.NotifyCloser, error) {
	fmt.Println("subsribe to on OnUpdate")

	return func() error {
		fmt.Println("unsubscribe to OnUpdate")
		return nil
	}, nil
}

package chipmonk

import (
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
)

type Friend struct {
	Name string
}

type Chipmunk struct {
	Friends map[string]*Friend
}

func Manage(c *Chipmunk) node.Node {
	// simple reflection works here to handle all management operations defined
	// in chipmunk.yang
	return &nodeutil.Node{Object: c}
}

package demo

import (
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
)

func manage(app *Contacts) node.Node {
	return nodeutil.ReflectChild(app)
}

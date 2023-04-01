package demo

import (
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
)

func manageJunkDrawer(app *JunkDrawer) node.Node {
	return nodeutil.ReflectChild(app)
}

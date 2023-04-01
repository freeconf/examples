package demo

import (
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

func TestManageJunk(t *testing.T) {
	app := &JunkDrawer{Info: make(map[string]interface{})}
	ypath := source.Dir(".")
	m := parser.RequireModule(ypath, "app-loose")
	b := node.NewBrowser(m, manageJunkDrawer(app))

	actual, err := nodeutil.WriteJSON(b.Root())
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, `{"size":0}`, actual)
}

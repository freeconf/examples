package demo

import (
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

func TestManageContact(t *testing.T) {
	app := &Contacts{}
	ypath := source.Dir(".")
	m := parser.RequireModule(ypath, "contacts")
	b := node.NewBrowser(m, manage(app))

	actual, err := nodeutil.WriteJSON(b.Root())
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, `{"size":0}`, actual)
}

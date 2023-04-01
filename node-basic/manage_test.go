package demo

import (
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

func TestManage(t *testing.T) {
	a := &App{}
	ypath := source.Dir(".")
	m := parser.RequireModule(ypath, "my-app")
	bwsr := node.NewBrowser(m, manage(a))

	root := bwsr.Root()
	actual, err := nodeutil.WriteJSON(root)
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, `{}`, actual)
}

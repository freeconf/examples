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
	b := &Bird{}
	ypath := source.Dir(".")
	m := parser.RequireModule(ypath, "bird")
	bwsr := node.NewBrowser(m, manage(b))

	root := bwsr.Root()
	actual, err := nodeutil.WriteJSON(root)
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, `{"location":"0,0"}`, actual)
}

package bear

import (
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

func TestManageReadOnly(t *testing.T) {
	bear := &Bear{Cubs: []*Cub{{Name: "bubbie"}}}
	ypath := source.Dir(".")
	m := parser.RequireModule(ypath, "bear-ro")
	b := node.NewBrowser(m, manageReadOnly(bear))

	actual, err := nodeutil.WriteJSON(b.Root())
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, `{"cub":[{"name":"bubbie"}]}`, actual)
}

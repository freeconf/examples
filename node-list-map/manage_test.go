package chipmonk

import (
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

func TestManage(t *testing.T) {
	cmunk := &Chipmunk{Friends: map[string]*Friend{"joe": {Name: "joe"}}}
	ypath := source.Dir(".")
	m := parser.RequireModule(ypath, "chipmunk")
	b := node.NewBrowser(m, Manage(cmunk))

	actual, err := nodeutil.WriteJSON(b.Root())
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, `{"friends":[{"name":"joe"}]}`, actual)
}

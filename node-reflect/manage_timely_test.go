package demo

import (
	"testing"
	"time"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

func TestManageTimely(t *testing.T) {
	app := &Timely{LastModified: time.Unix(99, 0)}
	ypath := source.Dir(".")
	m := parser.RequireModule(ypath, "timely")
	b := node.NewBrowser(m, manageTimely(app))

	actual, err := nodeutil.WriteJSON(b.Root())
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, `{"lastModified":99}`, actual)
}

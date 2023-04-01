package demo

import (
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
)

// Load a YANG module from a string!
var mstr = `
module mystery {
	namespace "freeconf.org";
	prefix "m";

	leaf name {
		type string;
	}
}
`

func TestMystery(t *testing.T) {

	m, err := parser.LoadModuleFromString(nil, mstr)
	fc.AssertEqual(t, nil, err)

	t.Run("struct", func(t *testing.T) {

		// create a node from an adhoc struct
		app := struct {
			Name string
		}{
			"john",
		}
		manage := nodeutil.ReflectChild(&app)

		// verify works
		b := node.NewBrowser(m, manage)
		actual, err := nodeutil.WriteJSON(b.Root())
		fc.AssertEqual(t, nil, err)
		fc.AssertEqual(t, `{"name":"john"}`, actual)
	})

	t.Run("map", func(t *testing.T) {

		// create a node from map
		app := map[string]interface{}{"name": "mary"}
		manage := nodeutil.ReflectChild(app)

		// verify works
		b := node.NewBrowser(m, manage)
		actual, err := nodeutil.WriteJSON(b.Root())
		fc.AssertEqual(t, nil, err)
		fc.AssertEqual(t, `{"name":"mary"}`, actual)
	})
}

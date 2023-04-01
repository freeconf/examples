package demo

import (
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

type Here struct {
	Penny int
}

func TestUnitTestingPartialYang(t *testing.T) {

	// where car.yang is stored
	ypath := source.Dir(".")

	// Define new YANG module on the fly that references the application
	// YANG file but we pull in just what we want
	m, err := parser.LoadModuleFromString(ypath, `
		module x {
			import deep {
				prefix "d";
			}

			// pull in just the piece we are interested in. Here it is
			// just a single penny
			uses d:here;
		}
	`)
	if err != nil {
		t.Fatal(err)
	}

	// We create a "browser" to just a unit of our application
	h := &Here{Penny: 1}
	manage := node.NewBrowser(m, nodeutil.ReflectChild(h))

	// Example : test getting config
	actual, err := nodeutil.WriteJSON(manage.Root())
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, `{"penny":1}`, actual)
}

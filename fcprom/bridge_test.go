package fcprom

import (
	"bytes"
	"flag"
	"strings"
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/source"

	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/val"
)

var updateFlag = flag.Bool("update", false, "update golden files instead of verifying against them")

const testModule = `module x {
	prefix "x";
	namespace "x";
	revision 0000-00-00;

	import metrics-extension {
		prefix "metrics";
	}

	leaf c {
		description "int32 counter";
		type int32;
		config false;
		metrics:counter;
	}

	leaf g {
		description "int32 gauge with
		  multi-line description
		";
		type int32;
		config false;
	}

	container y {
		config false;

		leaf z { 
			description "float gauge";
			type decimal64;
		}
	}

	leaf m {
		description "should not show as it is a configurable";
		type int32;
	}

	list f {
		config false;
		key "a b";
		metrics:multivariate;

		leaf a {
			type string;
		}

		leaf b {
			type string;
		}

		leaf g {			
			type int32;
		}
	}
}`

func TestBridge(t *testing.T) {
	ypath := source.Any(
		source.Named("x", strings.NewReader(testModule)),
		source.Dir("."))
	m, err := parser.LoadModuleFromString(ypath, testModule)
	if err != nil {
		t.Fatal(err)
	}
	n := &nodeutil.Basic{}
	n.OnField = func(r node.FieldRequest, hnd *node.ValueHandle) error {
		hnd.Val = val.Int32(99)
		return nil
	}
	n.OnChild = func(r node.ChildRequest) (node.Node, error) {
		return n, nil
	}
	n.OnNext = func(r node.ListRequest) (node.Node, []val.Value, error) {
		if r.Row >= 2 {
			return nil, nil, nil
		}
		keys, _ := node.NewValues(r.Meta.KeyMeta(), r.Row, r.Row*10)
		return n, keys, nil
	}
	bwsr := node.NewBrowser(m, n)
	x := newExporter()
	metricNode := x.node("x", []string{})
	root := bwsr.Root()
	defer root.Release()
	sel, err := root.Constrain("content=nonconfig")
	defer sel.Release()
	if err := sel.InsertInto(metricNode); err != nil {
		t.Fatal(err)
	}
	var actual bytes.Buffer
	writeMetrics(&actual, x.metrics)
	fc.Gold(t, *updateFlag, actual.Bytes(), "./gold/bridge1.txt")
}

func TestClean(t *testing.T) {
	fc.AssertEqual(t, "prom_bridge", metricName("prom-bridge"))
}

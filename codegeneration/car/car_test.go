package car

import (
	"fmt"

	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodes"
	"github.com/freeconf/yang/parser"
)

func Example_Car() {
	c := &Car{
		Speed: 10,
		Tire: []*Tire{
			&Tire{
				Size: "H15",
			},
		},
	}
	m := parser.RequireModule(meta.PathStreamSource(".."), "car")
	b := node.NewBrowser(m, CarNode(c))
	actual, err := nodes.WriteJSON(b.Root())
	if err != nil {
		panic(err)
	}
	// Output:
	// {"tire":[{"pos":0,"size":"H15","worn":false,"wear":0,"flat":false}],"speed":10,"miles":0,"lastRotation":0,"running":false}
	fmt.Println(actual)
}

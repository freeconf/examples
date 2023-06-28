package demo

import (
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/val"
)

func manage(b *Bird) node.Node {
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(b), // handles Name
		OnField: func(p node.Node, r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {
			case "location":
				hnd.Val = val.String(b.GetCoordinates())
			default:
				// delegates to ReflectChild for name
				return p.Field(r, hnd)
			}
			return nil
		},
	}
}

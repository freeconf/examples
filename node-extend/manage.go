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
				if r.Write {
					b.Location.Set(hnd.Val.String())
				} else {
					hnd.Val = val.String(b.Location.Get())
				}
			default:
				// delegates to ReflectChild
				return p.Field(r, hnd)
			}
			return nil
		},
	}
}

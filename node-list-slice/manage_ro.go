package bear

import (
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/val"
)

func manageReadOnly(b *Bear) node.Node {
	return &nodeutil.Basic{
		OnChild: func(r node.ChildRequest) (child node.Node, err error) {
			switch r.Meta.Ident() {
			case "cub":
				return manageCubsReadOnly(b), nil
			}
			return nil, nil
		},
	}
}

func manageCubsReadOnly(b *Bear) node.Node {
	return &nodeutil.Basic{
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			key := r.Key
			var found *Cub

			if key != nil {
				name := key[0].String()
				for _, cub := range b.Cubs {
					if cub.Name == name {
						found = cub
						break
					}
				}
			} else if r.Row < len(b.Cubs) {
				found = b.Cubs[r.Row]
				key = []val.Value{val.String(found.Name)}
			}
			if found != nil {
				return nodeutil.ReflectChild(found), key, nil
			}
			return nil, nil, nil
		},
	}
}

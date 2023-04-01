package bear

import (
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/val"
)

func manage(b *Bear) node.Node {
	return &nodeutil.Basic{
		OnChild: func(r node.ChildRequest) (child node.Node, err error) {
			switch r.Meta.Ident() {
			case "cub":
				return manageCubs(b), nil
			}
			return nil, nil
		},
	}
}

func manageCubs(b *Bear) node.Node {
	return &nodeutil.Basic{
		OnNextItem: func(r node.ListRequest) nodeutil.BasicNextItem {
			var found *Cub
			return nodeutil.BasicNextItem{
				New: func() error {
					name := r.Key[0].String()
					found = &Cub{Name: name}
					b.Cubs = append(b.Cubs, found)
					return nil
				},
				GetByKey: func() error {
					name := r.Key[0].String()
					for _, cub := range b.Cubs {
						if cub.Name == name {
							found = cub
							break
						}
					}
					return nil
				},
				DeleteByKey: func() error {
					name := r.Key[0].String()
					copy := make([]*Cub, 0, len(b.Cubs))
					for _, cub := range b.Cubs {
						if cub.Name != name {
							copy = append(copy, cub)
						}
					}
					b.Cubs = copy
					return nil
				},
				GetByRow: func() ([]val.Value, error) {
					var key []val.Value
					if r.Row < len(b.Cubs) {
						found = b.Cubs[r.Row]
						key = []val.Value{val.String(found.Name)}
					}
					return key, nil
				},
				Node: func() (node.Node, error) {
					if found != nil {
						return nodeutil.ReflectChild(found), nil
					}
					return nil, nil
				},
			}
		},
	}
}

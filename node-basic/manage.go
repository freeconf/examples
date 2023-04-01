package demo

import (
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
)

func manage(a *App) node.Node {
	return &nodeutil.Basic{
		OnChild: func(r node.ChildRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "users":
				return nodeutil.ReflectChild(a.users), nil
			case "fonts":
				return nodeutil.ReflectChild(a.fonts), nil
			case "bagels":
				return nodeutil.ReflectChild(a.bagels), nil
			}
			return nil, nil
		},
	}
}

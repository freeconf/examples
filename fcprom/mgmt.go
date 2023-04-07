package fcprom

import (
	"time"

	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/val"
)

func Manage(b *Bridge) node.Node {
	return &nodeutil.Basic{
		OnChild: func(r node.ChildRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "service":
				return serviceNode(b), nil
			case "modules":
				return nodeutil.ReflectChild(&b.Modules), nil
			case "render":
				return renderMetrics(b.RenderMetrics), nil
			}
			return nil, nil
		},
	}
}

func renderMetrics(m RenderMetrics) node.Node {
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(&m),
		OnField: func(p node.Node, r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {
			case "duration":
				hnd.Val = val.Int64(m.Duration / time.Millisecond)
			default:
				return p.Field(r, hnd)
			}
			return nil
		},
	}
}

func serviceNode(b *Bridge) node.Node {
	options := b.Options()
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(&options),
		OnEndEdit: func(p node.Node, r node.NodeRequest) error {
			if err := p.EndEdit(r); err != nil {
				return err
			}
			return b.Apply(options)
		},
	}
}

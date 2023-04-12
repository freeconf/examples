package fcinflux

import (
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/val"
)

// nodeWtr implements the FreeCONF node interface and acts like a "writer"
// when traversing a node tree with Selection.UpsertInto(<thisWriter>). There
// are other ways to accomplish walking a tree but this seemed easiest.
func nodeWtr(s sink, m Metric) node.Node {
	fields := make(map[string]interface{})
	return &nodeutil.Basic{
		OnField: func(r node.FieldRequest, h *node.ValueHandle) error {
			fields[r.Meta.Ident()] = h.Val.Value()
			return nil
		},
		OnChild: func(r node.ChildRequest) (node.Node, error) {
			return nodeWtr(s, m), nil
		},
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			return nodeWtr(s, m), r.Key, nil
		},
		OnEndEdit: func(r node.NodeRequest) error {
			if len(fields) > 0 {
				m.Name = r.Selection.Path.String()
				return s.send(r.Selection.Context, m, fields)
			}
			return nil
		},
	}
}

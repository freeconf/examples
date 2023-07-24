package car

import (
	"time"

	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/val"
)

/////////////////////////
// C A R    M A N A G E M E N T
//  Bridge from model to car app.

// Manage is root handler from car.yang. i.e. module car { ... }
func Manage(c *Car) node.Node {

	// Extend and Reflect form a powerful combination, we're letting reflect do a lot of the CRUD
	// when the yang file matches the field names.  But we extend reflection
	// to add as much custom behavior as we want
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(c),

		// implement navigation by containers and lists defined in yang file
		OnChild: func(p node.Node, r node.ChildRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "tire":
				return manageTires(c.Tire), nil
			default:
				// delegate back to "Base" handler should there be any default
				// handling
				return p.Child(r)
			}
		},

		// implement RPCs
		OnAction: func(p node.Node, r node.ActionRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "start":
				c.Start()
			case "stop":
				c.Stop()
			case "rotateTires":
				c.rotateTires()
			case "replaceTires":
				c.replaceTires()
			case "reset":
				c.Miles = 0
			}
			return nil, nil
		},

		// implement yang notifications which are really just events
		OnNotify: func(p node.Node, r node.NotifyRequest) (node.NotifyCloser, error) {
			switch r.Meta.Ident() {
			case "update":
				// can use an adhoc struct send event
				sub := c.OnUpdate(func(e updateEvent) {
					msg := struct {
						Event int
					}{
						Event: int(e),
					}
					// events are nodes too
					r.Send(nodeutil.ReflectChild(&msg))
				})

				// we return a close **function**, we are not actually closing here
				return sub.Close, nil
			}
			return nil, nil
		},

		// implement fields that are not automatically handled by reflection.
		OnField: func(p node.Node, r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {
			case "pollInterval":
				// another option for this is register field handler for reflection so that
				// it would handle this w/o having to implement OnField here
				if r.Write {
					c.PollInterval = time.Duration(hnd.Val.Value().(int)) * time.Millisecond
				} else {
					hnd.Val = val.Int32(int(c.PollInterval / time.Millisecond))
				}
			default:
				return p.Field(r, hnd)
			}
			return nil
		},
	}
}

// tiresNode handles list of tires.
//
//	list tire { ... }
func manageTires(tires []*Tire) node.Node {
	return &nodeutil.Basic{
		OnNextItem: func(r node.ListRequest) nodeutil.BasicNextItem {
			var t *Tire
			return nodeutil.BasicNextItem{
				GetByKey: func() error {
					// request for a specific tire by key (pos)
					pos := r.Key[0].Value().(int)
					if pos < len(tires) {
						t = tires[pos]
					}
					return nil
				},
				GetByRow: func() ([]val.Value, error) {
					// request for nth tire in list
					if r.Row < len(tires) {
						t = tires[r.Row]
						return []val.Value{val.Int32(r.Row)}, nil
					}
					return nil, nil
				},
				Node: func() (node.Node, error) {
					if t != nil {
						return ManageTire(t), nil
					}
					return nil, nil
				},
			}
		},
	}
}

// TireNode handles each tire node.  Everything *inside* list tire { ... }
func ManageTire(t *Tire) node.Node {
	// again, let reflection do a lot of the work with one extension to handle replace tire
	// action
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(t),
		OnAction: func(parent node.Node, r node.ActionRequest) (output node.Node, err error) {
			switch r.Meta.Ident() {
			case "replace":
				t.replace()
			}
			return nil, nil
		},
	}
}

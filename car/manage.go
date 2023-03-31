package car

import (
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/val"
)

// Useful (but not required) technique to organize management functions into
// a namespace. Generally all state is kept in app.
type manage struct{}

/////////////////////////
// C A R    M A N A G E M E N T
//  Bridge from model to car app.

// Manage is root handler from car.yang. i.e. module car { ... }
func Manage(c *Car) node.Node {
	var m manage

	// Powerful combination, we're letting reflect do a lot of the CRUD
	// when the yang file matches the field names.  But we extend reflection
	// to add as much custom behavior as we want
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(c),

		// drilling into child objects defined by yang file
		OnChild: func(p node.Node, r node.ChildRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "tire":
				return m.tires(c.Tire), nil
			default:
				// delegate back to "Base" handler should there be any default
				// handling
				return p.Child(r)
			}
		},

		// RPCs
		OnAction: func(p node.Node, r node.ActionRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "rotateTires":
				c.rotateTires()
			case "replaceTires":
				c.replaceTires()
			case "reset":
				c.Miles = 0
			}
			return nil, nil
		},

		// Events
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
					r.Send(nodeutil.ReflectChild(&msg))
				})

				// we return a close **function**, we are not actually closing here
				return sub.Close, nil
			}
			return nil, nil
		},

		// override OnEndEdit just to know when car has been created and
		// fully initialized so we can start the car running
		OnEndEdit: func(p node.Node, r node.NodeRequest) error {
			// allow reflection node handler to finish, this is where defaults
			// get set.
			if err := p.EndEdit(r); err != nil {
				return err
			}
			c.Start()
			return nil
		},
	}
}

// tiresNode handles list of tires.
//
//	list tire { ... }
func (m manage) tires(tires []*Tire) node.Node {
	return &nodeutil.Basic{
		OnNextItem: func(r node.ListRequest) nodeutil.BasicNextItem {
			var t *Tire
			return nodeutil.BasicNextItem{
				GetByKey: func() error {
					pos := r.Key[0].Value().(int)
					if pos < len(tires) {
						t = tires[pos]
					}
					return nil
				},
				GetByRow: func() ([]val.Value, error) {
					// request for nth item in list
					if r.Row < len(tires) {
						t = tires[r.Row]
						return []val.Value{val.Int32(r.Row)}, nil
					}
					return nil, nil
				},
				Node: func() (node.Node, error) {
					if t != nil {
						// let reflection do a lot of the work
						return m.tire(t), nil
					}
					return nil, nil
				},
			}
		},
	}
}

// TireNode handles each tire node.  Everything *inside* list tire { ...}
func (m manage) tire(t *Tire) node.Node {
	// let reflection do a lot of the work
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

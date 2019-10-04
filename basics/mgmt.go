package basics

import (
	"errors"
	"fmt"

	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodes"
	"github.com/freeconf/yang/val"
)

// Manage is a bridge from model to the car application.  This is the only place where you
// couple your application code to FreeCONF.
func Manage(c *Car) node.Node {

	// Powerful combination, we're letting reflect do a lot of the CRUD
	// when the yang file matches the field names.  But we extend reflection
	// to add as much custom behavior as we want
	return &nodes.Extend{

		// Reflection
		Base: nodes.ReflectChild(c),

		// CRUD drilling into child objects defined by yang file
		OnChild: func(p node.Node, r node.ChildRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "tire":
				return tiresNode(c), nil
			case "specs":
				// knows how to r/w config from a map
				return nodes.ReflectChild(c.Specs), nil
			default:
				// return control back to handler we're extending, in this case
				// it's reflection
				return p.Child(r)
			}
			return nil, nil
		},

		// Functions
		OnAction: func(p node.Node, r node.ActionRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "rotateTires":
				c.rotateTires()
			case "replaceTires":
				c.replaceTires()
			}
			return nil, nil
		},

		// Events
		OnNotify: func(p node.Node, r node.NotifyRequest) (node.NotifyCloser, error) {
			switch r.Meta.Ident() {
			case "update":

				// typical application event listener code bridging to FreeCONF
				sub := c.OnUpdate(func(*Car) {

					// cleverly reuses node handler to send car data
					r.Send(Manage(c))

				})

				// NOTE: we return a close function, we are not actually closing
				// here
				return func() error {
					c.UnsubscribeOnUpdate(sub)
					return nil
				}, nil
			}
			return nil, nil
		},

		// override OnEndEdit just to just to know when car has been created and
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
//     list tire { ... }
func tiresNode(c *Car) node.Node {
	return &nodes.Basic{
		// Handling lists are
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			var t *tire
			key := r.Key

			// request for specific item in list
			if key != nil {
				if r.New {
					nextPos := len(c.Tire)
					if key[0].Value().(int) != nextPos {
						msg := fmt.Sprintf("pos must be next sequential value %d", nextPos)
						return nil, nil, errors.New(msg)
					}
					t = &tire{
						Pos: nextPos,
					}
					c.Tire = append(c.Tire, t)
				} else {
					pos := key[0].Value().(int)
					if pos >= len(c.Tire) {
						return nil, nil, nil
					}
					t = c.Tire[pos]
				}
			} else {
				// request for nth item in list
				if r.Row < len(c.Tire) {
					t = c.Tire[r.Row]
					key = []val.Value{val.Int32(r.Row)}
				}
			}
			if t != nil {
				return tireNode(t), key, nil
			}
			return nil, nil, nil
		},
	}
}

// tireNode handles each tire node.  Everything *inside* list tire { ...}
func tireNode(t *tire) node.Node {

	// Again, let reflection do a lot of the work
	return &nodes.Extend{
		Base: nodes.ReflectChild(t),

		OnField: func(p node.Node, r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {

			case "worn":
				// worn is a method call, so our current reflection handler doesn't
				// check for that.  Maybe you reflection handler would.
				hnd.Val = val.Bool(t.Worn())

			default:
				return p.Field(r, hnd)
			}
			return nil
		},
	}
}

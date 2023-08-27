package car

import (
	"reflect"
	"time"

	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
)

// ///////////////////////
// C A R    M A N A G E M E N T
//
// Manage your car application using FreeCONF library according to the car.yang
// model file.
//
// Manage is root handler from car.yang. i.e. module car { ... }
func Manage(car *Car) node.Node {

	// We're letting reflect do a lot of the work when the yang file matches
	// the field names and methods in the objects.  But we extend reflection
	// to add as much custom behavior as we want
	return &nodeutil.Node{

		// Initial object. Note: as the tree is traversed, new Node instances
		// will have different values in their Object reference
		Object: car,

		// implement RPCs
		OnAction: func(n *nodeutil.Node, r node.ActionRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "reset":
				// here we implement custom handler for action just as an example
				// If there was a Reset() method on car the default OnAction
				// handler would be able to call Reset like all the other functions
				car.Miles = 0
			default:
				// all the actions like start, stop and rotateTire are called
				// thru reflecton here because their method names align with
				// the YANG.
				return n.DoAction(r)
			}
			return nil, nil
		},

		// implement yang notifications (which are really just event streams)
		OnNotify: func(p *nodeutil.Node, r node.NotifyRequest) (node.NotifyCloser, error) {
			switch r.Meta.Ident() {
			case "update":
				// can use an adhoc struct send event
				sub := car.OnUpdate(func(e UpdateEvent) {
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

			// there is no default implementation at this time, all notification streams
			// require custom handlers.
			return p.Notify(r)
		},

		// implement fields that are not automatically handled by reflection.
		OnRead: func(p *nodeutil.Node, r meta.Definition, t reflect.Type, v reflect.Value) (reflect.Value, error) {
			if l, ok := r.(meta.Leafable); ok {
				// use "units" in yang to tell what to convert.
				//
				// Other useful ways to intercept custom reflection reads:
				// 1.) incoming reflect.Type t
				// 2.) field name used in yang or some pattern of the name (suffix, prefix, regex)
				// 3.) yang extension
				// 4.) any combination of these
				if l.Units() == "millisecs" {
					return reflect.ValueOf(v.Int() / int64(time.Millisecond)), nil
				}
			}
			return v, nil
		},

		// Generally the reverse of what is handled in OnRead
		OnWrite: func(p *nodeutil.Node, r meta.Definition, t reflect.Type, v reflect.Value) (reflect.Value, error) {
			if l, ok := r.(meta.Leafable); ok {
				if l.Units() == "millisecs" {
					d := time.Duration(v.Int()) * time.Millisecond
					return reflect.ValueOf(d), nil
				}
			}
			return v, nil
		},
	}
}

package fcinflux

import (
	"fmt"
	"time"

	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/val"
)

// Manage implements the API to management protocols like RESTCONF or gNMI to this fc-infludb
// module itself.
func Manage(svc *Service) node.Node {
	return &nodeutil.Basic{
		OnChild: func(r node.ChildRequest) (child node.Node, err error) {
			switch r.Meta.Ident() {
			case "options":
				return manageOptions(svc), nil
			}
			return nil, nil
		},
	}
}

// manageOptions gathers the new configuration options then applies them at the end
func manageOptions(svc *Service) node.Node {
	o := svc.Options()
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(&o),
		OnField: func(parent node.Node, r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {
			case "frequency":
				if r.Write {
					secs := hnd.Val.Value().(uint64)
					o.Frequency = time.Duration(secs) * time.Second
				} else {
					secs := o.Frequency / time.Second
					hnd.Val = val.UInt64(secs)
				}
			case "tags":
				if r.Write {
					// reflect returns a map[string]interface{} and so we need to
					// custom handle the setter
					tags := hnd.Val.Value().(map[string]interface{})
					o.Tags = decodeTags(tags)
				} else {
					// reflect can convert from map[string]string to anything w/o issue tho
					hnd.Val = val.Any{Thing: o.Tags}
				}
			default:
				return parent.Field(r, hnd)
			}
			return nil
		},
		OnEndEdit: func(p node.Node, r node.NodeRequest) error {
			if err := p.EndEdit(r); err != nil {
				return err
			}
			return svc.ApplyOptions(o)
		},
	}
}

func decodeTags(in map[string]interface{}) map[string]string {
	out := make(map[string]string, len(in))
	for key, v := range in {
		// incase numbers were given, convert them to strings. could probably
		// check for containers and lists there but instead we just smash them
		// into strings.
		out[key] = fmt.Sprintf("%s", v)
	}
	return out
}

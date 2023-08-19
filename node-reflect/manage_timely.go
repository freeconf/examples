package demo

import (
	"reflect"
	"time"

	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
)

func manageTimely(t *Timely) node.Node {
	timeType := reflect.TypeOf(time.Time{})
	return &nodeutil.Node{
		Object: t,
		OnRead: func(ref *nodeutil.Node, m meta.Definition, t reflect.Type, v reflect.Value) (reflect.Value, error) {
			switch t {
			case timeType:
				at := v.Interface().(time.Time)
				return reflect.ValueOf(at.Unix()), nil
			}
			return v, nil
		},
		OnWrite: func(ref *nodeutil.Node, m meta.Definition, t reflect.Type, v reflect.Value) (reflect.Value, error) {
			switch t {
			case timeType:
				return reflect.ValueOf(time.Unix(v.Int(), 0)), nil
			}
			return v, nil
		},
	}
}

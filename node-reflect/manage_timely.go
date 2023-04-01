package demo

import (
	"reflect"
	"time"

	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/val"
)

var timeHandler = nodeutil.ReflectField{
	When: nodeutil.ReflectFieldByType(reflect.TypeOf(time.Time{})),
	OnRead: func(leaf meta.Leafable, fieldname string, elem reflect.Value, fieldElem reflect.Value) (val.Value, error) {
		t := fieldElem.Interface().(time.Time)
		return val.Int64(t.Unix()), nil
	},
	OnWrite: func(leaf meta.Leafable, fieldname string, elem reflect.Value, fieldElem reflect.Value, v val.Value) error {
		t := time.Unix(v.Value().(int64), 0)
		fieldElem.Set(reflect.ValueOf(t))
		return nil
	},
}

func manageTimely(t *Timely) node.Node {
	return nodeutil.Reflect{
		OnField: []nodeutil.ReflectField{
			timeHandler,
		},
	}.Object(t)
}

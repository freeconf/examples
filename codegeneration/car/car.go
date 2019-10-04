package car

// This file was auto-generated. Do not edit

import (
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodes"
	"github.com/freeconf/yang/val"
)

/* Car - Vehicle of sorts
 */
type Car struct {
	Speed        int
	Miles        int64
	LastRotation int64
	Running      bool
	Engine       *Engine
	Tire         []*Tire
}

/* Engine - powers the car
 */
type Engine struct {
	Specs *Specs
}

/* Specs - details about the car
 */
type Specs struct {
	Horsepower int
}

/* Tire - rubber circular part that makes contact with road
 */
type Tire struct {
	Pos  int
	Size string
	Worn bool
	Wear float64
	Flat bool
}

/*  CarNode - Management of Vehicle of sorts
 */
func CarNode(o *Car) node.Node {
	return &nodes.Basic{
		OnField: func(r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {
			case "speed":
				hnd.Val = val.Int32(o.Speed)
			case "miles":
				hnd.Val = val.Int64(o.Miles)
			case "lastRotation":
				hnd.Val = val.Int64(o.LastRotation)
			case "running":
				hnd.Val = val.Bool(o.Running)
			}
			return nil
		},
		OnChild: func(r node.ChildRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "engine":
				if o.Engine != nil {
					return EngineNode(o.Engine), nil
				}
			case "tire":
				if len(o.Tire) > 0 {
					return TireListNode(o.Tire), nil
				}
			}
			return nil, nil
		},
		OnAction: func(r node.ActionRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "start":
				return o.doStart(r)
			}
			return nil, nil
		},
		OnNotify: func(r node.NotifyRequest) (node.NotifyCloser, error) {
			switch r.Meta.Ident() {
			case "update":
				return o.onUpdate(r)
			}
			return nil, nil
		},
	}
}

func TireListNode(o []*Tire) node.Node {
	return &nodes.Basic{
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			if r.Row < len(o) {
				return TireNode(o[r.Row]), nil, nil
			}
			return nil, nil, nil
		},
	}
}

/*  EngineNode - Management of powers the car
 */
func EngineNode(o *Engine) node.Node {
	return &nodes.Basic{
		OnChild: func(r node.ChildRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "specs":
				if o.Specs != nil {
					return SpecsNode(o.Specs), nil
				}
			}
			return nil, nil
		},
	}
}

/*  SpecsNode - Management of details about the car
 */
func SpecsNode(o *Specs) node.Node {
	return &nodes.Basic{
		OnField: func(r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {
			case "horsepower":
				hnd.Val = val.Int32(o.Horsepower)
			}
			return nil
		},
	}
}

/*  TireNode - Management of rubber circular part that makes contact with road
 */
func TireNode(o *Tire) node.Node {
	return &nodes.Basic{
		OnField: func(r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {
			case "pos":
				hnd.Val = val.Int32(o.Pos)
			case "size":
				hnd.Val = val.String(o.Size)
			case "worn":
				hnd.Val = val.Bool(o.Worn)
			case "wear":
				hnd.Val = val.Decimal64(o.Wear)
			case "flat":
				hnd.Val = val.Bool(o.Flat)
			}
			return nil
		},
	}
}

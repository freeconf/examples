package fcslack

import (
	"errors"

	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/val"
)

func Manage(c *Service) node.Node {
	return &nodeutil.Basic{
		OnChild: func(r node.ChildRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "subscription":
				return manageNotifications(c), nil
			case "client":
				return manageSlackClient(c.slack), nil
			}
			return nil, nil
		},
		OnNotify: func(r node.NotifyRequest) (node.NotifyCloser, error) {
			switch r.Meta.Ident() {
			case "errors":
				sub := c.OnError(func(err error, sub *subscription) {
					msg := map[string]interface{}{
						"module": sub.Module,
						"path":   sub.Path,
						"error":  err.Error(),
					}
					r.Send(nodeutil.ReflectChild(msg))
				})
				return sub.Close, nil
			}
			return nil, nil
		},
	}
}

func manageSlackClient(client *slack) node.Node {
	opts := client.options()
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(&opts),
		OnEndEdit: func(parent node.Node, r node.NodeRequest) error {
			if err := parent.EndEdit(r); err != nil {
				return err
			}
			return client.apply(opts)
		},
	}
}

var errNoKey = errors.New("must specify module and path")

func manageNotifications(s *Service) node.Node {
	index := node.NewIndex(s.notifications)
	return &nodeutil.Basic{
		OnNextItem: func(r node.ListRequest) nodeutil.BasicNextItem {
			var found *subscription
			var key string
			if len(r.Key) == 2 {
				key = subscriptionKey(r.Key[0].String(), r.Key[1].String())
			}
			return nodeutil.BasicNextItem{
				GetByKey: func() error {
					found = s.notifications[key]
					return nil
				},
				GetByRow: func() ([]val.Value, error) {
					if r.Row >= index.Len() {
						return nil, nil
					}
					v := index.NextKey(r.Row)
					if v == node.NO_VALUE {
						return nil, nil
					}
					id := v.String()
					found = s.notifications[id]
					return []val.Value{val.String(found.Module), val.String(found.Path)}, nil
				},
				New: func() error {
					if key == "" {
						return errNoKey
					}
					found = &subscription{Module: r.Key[0].String(), Path: r.Key[1].String()}
					return nil
				},
				DeleteByKey: func() error {
					s.remove(key)
					return nil
				},
				Node: func() (node.Node, error) {
					if found != nil {
						return manageSub(s, found), nil
					}
					return nil, nil
				},
			}
		},
	}
}

func manageSub(s *Service, e *subscription) node.Node {
	return &nodeutil.Extend{
		Base: nodeutil.ReflectChild(e),
		OnField: func(parent node.Node, r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {
			case "active":
				hnd.Val = val.Bool(e.Closer != nil)
			default:
				return parent.Field(r, hnd)
			}
			return nil
		},
		OnEndEdit: func(parent node.Node, r node.NodeRequest) error {
			if err := parent.EndEdit(r); err != nil {
				return err
			}
			if !r.Delete {
				return s.add(e)
			}
			return nil
		},
	}
}

package car

import (
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

// NOTE: you don't really want to repeat what is already in car_test.go, just want to
// exercise the custom code in manager_test.go
func TestManage(t *testing.T) {
	c, mgr := setupNewTestManager()
	root := mgr.Root()

	// Read all config
	actual, err := nodeutil.WriteJSON(root.Find("?content=config"))
	fc.AssertEqual(t, nil, err)
	expected := `{"speed":1000,"tire":[{"pos":0,"size":"15"},{"pos":1,"size":"15"},{"pos":2,"size":"15"},{"pos":3,"size":"15"}]}`
	fc.AssertEqual(t, expected, actual)

	// access car and verify w/API
	fc.AssertEqual(t, false, c.Running)

	// setup events stream reader
	events := make(chan string)
	unsub, err := root.Find("update").Notifications(func(n node.Notification) {
		event, _ := nodeutil.WriteJSON(n.Event)
		events <- event
	})
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, 1, c.listeners.Len())

	// write config starts car
	err = root.UpdateFrom(nodeutil.ReadJSON(`{"speed":1000}`)).LastErr
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, 1000, c.Speed)

	// should be first event
	fc.AssertEqual(t, `{"event":"carStarted"}`, <-events)
	fc.AssertEqual(t, true, c.Running)

	// unsubscribe
	unsub()
	fc.AssertEqual(t, 0, c.listeners.Len())

	// hit all the RPCs
	fc.AssertEqual(t, nil, root.Find("rotateTires").Action(nil).LastErr)
	fc.AssertEqual(t, nil, root.Find("replaceTires").Action(nil).LastErr)
	fc.AssertEqual(t, nil, root.Find("reset").Action(nil).LastErr)
	fc.AssertEqual(t, nil, root.Find("tire=0/replace").Action(nil).LastErr)
}

// no server, just your app and management API.  Testing
// just part of the managment API is also possible but here
// we create a whole car
func setupNewTestManager() (*Car, *node.Browser) {
	c := New()
	ypath := source.Path("../yang")
	m := parser.RequireModule(ypath, "car")
	return c, node.NewBrowser(m, Manage(c))
}

package car

import (
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

// Test the car management logic in manage.go
func TestManage(t *testing.T) {

	// setup
	ypath := source.Path("../yang")
	mod := parser.RequireModule(ypath, "car")
	app := New()

	// no web server needed, just your app and management function.
	brwsr := node.NewBrowser(mod, Manage(app))
	root := brwsr.Root()

	// read all config
	currCfg, err := nodeutil.WriteJSON(find(root, "?content=config"))
	fc.AssertEqual(t, nil, err)
	expected := `{"speed":1000,"pollInterval":1000,"tire":[{"pos":0,"size":"H15"},{"pos":1,"size":"H15"},{"pos":2,"size":"H15"},{"pos":3,"size":"H15"}]}`
	fc.AssertEqual(t, expected, currCfg)

	// access car and verify w/API
	fc.AssertEqual(t, false, app.Running)

	// setup event listener, verify events later
	events := make(chan string)
	unsub, err := find(root, "update").Notifications(func(n node.Notification) {
		event, _ := nodeutil.WriteJSON(n.Event)
		events <- event
	})
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, 1, app.listeners.Len())

	// write config starts car
	err = root.UpdateFrom(nodeutil.ReadJSON(`{"speed":2000}`))
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, 2000, app.Speed)

	// start car
	fc.AssertEqual(t, nil, justErr(find(root, "start").Action(nil)))

	// should be first event
	fc.AssertEqual(t, `{"event":"carStarted"}`, <-events)
	fc.AssertEqual(t, true, app.Running)

	// unsubscribe
	unsub()
	fc.AssertEqual(t, 0, app.listeners.Len())

	// hit all the RPCs
	fc.AssertEqual(t, nil, justErr(find(root, "rotateTires").Action(nil)))
	fc.AssertEqual(t, nil, justErr(find(root, "replaceTires").Action(nil)))
	fc.AssertEqual(t, nil, justErr(find(root, "reset").Action(nil)))
	fc.AssertEqual(t, nil, justErr(find(root, "tire=0/replace").Action(nil)))
}

func find(sel *node.Selection, path string) *node.Selection {
	found, err := sel.Find(path)
	if err != nil {
		panic(err)
	}
	return found
}

func justErr(_ *node.Selection, err error) error {
	return err
}

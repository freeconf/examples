package basics_test

import (
	"strings"
	"testing"

	"github.com/freeconf/examples/basics"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

func TestManagement(t *testing.T) {
	yangPath := source.Dir(".")

	// parse YANG
	m := parser.RequireModule(yangPath, "car")

	// attach YANG to management code to create management browser object
	// use browser object to test interface.  Notice that http server is unnec.
	b := node.NewBrowser(m, basics.Manage(basics.New()))

	// Configuration
	if err := b.Root().Set("speed", 1); err != nil {
		t.Error(err)
	}

	// Operations
	if err := b.Root().Find("replaceTires").Action(nil).LastErr; err != nil {
		t.Error(err)
	}

	// Alerts
	wait := make(chan struct{})
	_, err := b.Root().Find("update").Notifications(func(update node.Notification) {
		tires, err := nodeutil.WriteJSON(update.Event.Find("tire?fields=worn%3Bflat"))
		if err != nil {
			t.Error(err)
		}

		// most likely worn, but could also be flat
		if strings.Contains(tires, "true") {
			wait <- struct{}{}
		}
	})
	if err != nil {
		t.Error(err)
	}
	t.Log("waiting for wear/flat...")
	<-wait

	// Metrics
	if v, err := b.Root().GetValue("running"); err != nil {
		t.Error(err)
	} else {
		if v.Value().(bool) != true {
			t.Error("not running")
		}
	}
}

package fcslack

import (
	"testing"

	"github.com/freeconf/restconf/device"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"

	"github.com/freeconf/yang/fc"
)

func TestManage(t *testing.T) {
	ypath := source.Dir(".")
	send := make(chan string)
	testBrowser := loadTestBrowser(send)
	d := device.New(ypath)
	d.AddBrowser(testBrowser)

	c := NewService(d)
	recv := make(chan msg)
	c.sink = func(m msg) error {
		recv <- m
		return nil
	}
	m := parser.RequireModule(ypath, "fc-slack")
	b := node.NewBrowser(m, Manage(c))
	n, err := nodeutil.ReadJSON(`{
		"subscription" : [{
			"module" : "m",
			"path" : "n"
		}]
	}`)
	fc.AssertEqual(t, nil, err)
	err = b.Root().UpsertFrom(n)
	fc.RequireEqual(t, nil, err)
	fc.RequireEqual(t, 1, len(c.notifications))

	send <- `{"x":"hi"}`
	actual := <-recv
	fc.AssertEqual(t, uint32(1), c.notifications["m:n"].Counter)
	fc.AssertEqual(t, `{"x":"hi"}`, actual.Text)

	actualStr, err := nodeutil.WriteJSON(b.Root())
	fc.RequireEqual(t, nil, err)
	fc.AssertEqual(t, `{"client":{"debug":false},"subscription":[{"module":"m","path":"n","counter":1,"active":true}]}`, actualStr)
}

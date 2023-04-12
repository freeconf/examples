package fcslack

import (
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
)

func TestSubscribe(t *testing.T) {
	send := make(chan string)
	b := loadTestBrowser(send)
	recv := make(chan msg)
	sink := func(m msg) error {
		recv <- m
		return nil
	}
	sub := &subscription{
		Module:  "m",
		Channel: "c",
		Path:    "n",
	}
	err := sub.subscribe(b, sink, func(err error, e *subscription) {
		panic(err)
	})
	fc.AssertEqual(t, nil, err)
	expected := `{"x":"hi"}`
	send <- expected
	actual := <-recv
	fc.AssertEqual(t, "c", actual.Channel)
	fc.AssertEqual(t, expected, actual.Text)
	sub.Close()
}

// create a sample browser that has single notification item
func loadTestBrowser(send <-chan string) *node.Browser {
	mstr := `
		module m {
			revision 0;
			notification n {
				leaf x {
					type string;
				}
			}
		}`
	m, err := parser.LoadModuleFromString(nil, mstr)
	if err != nil {
		panic(err)
	}
	n := &nodeutil.Basic{
		OnNotify: func(r node.NotifyRequest) (node.NotifyCloser, error) {
			go func() {
				r.Send(nodeutil.ReadJSON(<-send))
			}()
			nop := func() error {
				return nil
			}
			return nop, nil
		},
	}
	return node.NewBrowser(m, n)
}

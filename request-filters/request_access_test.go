package demo

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/freeconf/restconf"
	"github.com/freeconf/restconf/device"
	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
	"github.com/freeconf/yang/val"
)

type App struct {
}

var module = `module x {
	namespace "freeconf.org";
	prefix "x";

	leaf messageFromRequest {
		config false;
		type string;
	}
}
`

type testContextKey string

const contextKey = testContextKey("requestKey")

func manageApp(a *App) node.Node {
	return &nodeutil.Basic{
		OnField: func(r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {
			case "messageFromRequest":
				msg := r.Selection.Context.Value(contextKey).(string)
				hnd.Val = val.String(msg)
			}
			return nil
		},
	}
}

func startServer() *restconf.Server {
	ypath := source.Dir("../yang")
	m, err := parser.LoadModuleFromString(ypath, module)
	if err != nil {
		panic(err)
	}
	app := &App{}
	b := node.NewBrowser(m, manageApp(app))
	d := device.New(ypath)
	d.AddBrowser(b)
	s := restconf.NewServer(d)

	// Register a filter to peek at request.  From here you can look at:
	//  * headers
	//  * certs
	//  * url params
	grabHeader := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
		msg := r.Header.Get("X-MESSAGE")

		// stuff your data into the context and it will be available to all node
		// navigation associated with this request.
		ctx = context.WithValue(ctx, contextKey, msg)
		return ctx, nil
	}
	s.Filters = append(s.Filters, grabHeader)

	d.ApplyStartupConfig(strings.NewReader(`{
		"fc-restconf" : {
			"web" : {
				"port" : ":9999"
			}
		}
	}
	`))

	return s
}

func TestRequestAccess(t *testing.T) {
	startServer()

	t.Run("request", func(t *testing.T) {

		req, err := http.NewRequest("GET", "http://localhost:9999/restconf/data/x:", nil)
		fc.AssertEqual(t, nil, err)
		req.Header.Add("X-MESSAGE", "hi")
		c := &http.Client{}
		resp, err := c.Do(req)
		fc.AssertEqual(t, nil, err)
		fc.AssertEqual(t, 200, resp.StatusCode)
		actual, err := ioutil.ReadAll(resp.Body)
		fc.AssertEqual(t, nil, err)
		fc.AssertEqual(t, `{"x:messageFromRequest":"hi"}`, string(actual))
	})
}

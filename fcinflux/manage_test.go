package fcinflux

import (
	"testing"
	"time"

	"github.com/freeconf/restconf/device"
	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/source"
)

func TestManage(t *testing.T) {
	ypath := source.Path(".")
	d := device.New(ypath)
	s := NewService(d)
	fc.AssertEqual(t, nil, d.Add("fc-influx", Manage(s)))
	b, err := d.Browser("fc-influx")
	fc.AssertEqual(t, nil, err)
	n, err := nodeutil.ReadJSON(`{
		"options" : {
			"connection": { 
				"addr": "http://localhost:8086/metrics",
				"apiToken" : "abc"
			},
			"organization" : "example.org",
			"bucket" : "bbb",
			"frequency": 10,
			"database": "demo",
			"tags": {
				"abc" : "123"
			}
		}
	}`)
	fc.AssertEqual(t, nil, err)
	err = b.Root().UpsertFrom(n)
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, 10*time.Second, s.options.Frequency)
	fc.AssertEqual(t, "abc", s.options.Connection.ApiToken)
	fc.AssertEqual(t, map[string]string{"abc": "123"}, s.options.Tags)
}

package fcprom

import (
	"bytes"
	"testing"

	"github.com/freeconf/restconf/device"
	"github.com/freeconf/yang/source"
)

func TestBridgeMgmt(t *testing.T) {
	ypath := source.Path("../yang:.")
	d := device.New(ypath)
	b := NewBridge(d)
	if err := d.Add("fc-prom", Manage(b)); err != nil {
		t.Fatal(err)
	}
	var actual bytes.Buffer
	if err := b.generate(&actual); err != nil {
		t.Fatal(err)
	}
	t.Log(actual.String())
}

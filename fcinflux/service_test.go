package fcinflux

import (
	"context"
	"testing"
	"time"

	"github.com/freeconf/examples/car"
	"github.com/freeconf/restconf/device"
	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/source"
)

func TestStart(t *testing.T) {
	ypath := source.Dir("../yang")
	d := device.New(ypath)
	c := car.New()
	fc.AssertEqual(t, nil, d.Add("car", car.Manage(c)))

	dummy := &dummyDriver{
		done: make(chan bool),
	}
	s := newService(d, dummy.driver)
	err := s.ApplyOptions(Options{
		Database:  "x",
		Frequency: time.Millisecond,
	})
	fc.AssertEqual(t, nil, err)
	<-dummy.done
	fc.AssertEqual(t, 5, dummy.counter)
}

type dummyDriver struct {
	counter int
	done    chan bool
}

func (d *dummyDriver) send(ctx context.Context, m Metric, fields map[string]interface{}) error {
	d.counter++
	return nil
}

func (d *dummyDriver) close(context.Context) error {
	d.done <- true
	return nil
}

func (d *dummyDriver) driver(Options) (sink, error) {
	return d, nil
}

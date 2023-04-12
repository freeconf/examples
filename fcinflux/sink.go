package fcinflux

import (
	"context"

	"github.com/freeconf/yang/fc"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// sink is interface between nodeWtr and influx service.  This allows of unit
// testing of each w/o the other
type sink interface {
	send(ctx context.Context, m Metric, fields map[string]interface{}) error
	close(ctx context.Context) error
}

// influxSink is concrete implemetion of sink that actually write to influx service
type influxSink struct {
	client influxdb2.Client
	api    api.WriteAPIBlocking
}

// this allows influx testing the Service with a dummy driver that returns dummy
// influx sink.  It helps when testing manage.go
type sinkDriver func(opts Options) (sink, error)

// newInfluxSink connects to InfluxDB service.  This is meant to be created for each
// batch write to influxdb
func newInfluxSink(opts Options) (sink, error) {
	conn := influxdb2.NewClient(opts.Connection.Addr, opts.Connection.ApiToken)
	api := conn.WriteAPIBlocking(opts.Organization, opts.Bucket)
	return &influxSink{client: conn, api: api}, nil
}

// close and flush data to influxdb
func (sink *influxSink) close(ctx context.Context) error {
	defer sink.client.Close()
	return sink.api.Flush(ctx)
}

// send data to influxdb
func (sink *influxSink) send(ctx context.Context, m Metric, fields map[string]interface{}) error {
	if fc.DebugLogEnabled() {
		// only build string if debugging is on so as to not waste string interpolation of there
		// isn't going to be a write
		fc.Debug.Printf("sending %s -> %v", m.Name, fields)
	}
	p := influxdb2.NewPoint(m.Name, m.Tags, fields, m.Time)
	return sink.api.WritePoint(ctx, p)
}

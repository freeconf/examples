package fcinflux

import (
	"context"
	"errors"
	"time"

	"github.com/freeconf/restconf/device"
	"github.com/freeconf/yang/fc"
)

// Service orchestrates all uploads to influxdb
type Service struct {
	device  device.Device
	options Options
	driver  sinkDriver
	ticker  *time.Ticker
}

// ConnectionOptions captures the InfluxDB connection details.  There more options to add
// here including TLS and other possible login types
type ConnectionOptions struct {
	Addr     string
	ApiToken string
}

// Config of this service
type Options struct {
	Connection    ConnectionOptions
	Database      string // added to every metric to distinguish metrics from multiple sources
	Organization  string
	Bucket        string
	Tags          map[string]string // default set of tags for all metrics
	Frequency     time.Duration     // how often to grab metrics
	IgnoreModules []string
	IgnorePaths   []string
}

// Metric is converted to a Influx WritePoint
type Metric struct {
	Name string
	Tags map[string]string
	Time time.Time
}

// NewService is constructor
func NewService(device device.Device) *Service {
	return newService(device, newInfluxSink)
}

// newService is constructor for unit test
func newService(device device.Device, driver sinkDriver) *Service {
	sink := &Service{
		device: device,
		driver: driver,
	}
	return sink
}

// Options is access to config by returning a copy we ensure they cannot be
// changed w/o calling Apply
func (svc *Service) Options() Options {
	return svc.options
}

var errNonZeroFrequency = errors.New("only non-zero frequency allowed")

// ApplyOptions will allow for changing config w/o restart.
func (svc *Service) ApplyOptions(options Options) error {
	if options.Frequency == 0 {
		return errNonZeroFrequency
	}
	svc.options = options
	if svc.ticker != nil {
		svc.ticker.Stop()
	}
	svc.ticker = time.NewTicker(svc.options.Frequency)
	go svc.Start(svc.ticker.C)
	return nil
}

// handleErr is for errors that happen in background only. It's possible
// this might want option to convert to notification or fail service
func (svc *Service) handleErr(err error) {
	fc.Err.Print(err)
}

// Start the polling of metrics.
func (svc *Service) Start(ticker <-chan time.Time) {
	for t := range ticker {
		ctx := context.Background()
		if t.IsZero() {
			break // ticker stopped
		}
		sink, err := svc.driver(svc.options)
		if err != nil {
			svc.handleErr(err)
			continue
		}
		if err = svc.run(ctx, sink); err != nil {
			svc.handleErr(err)
			continue
		}
		if err = sink.close(ctx); err != nil {
			svc.handleErr(err)
			continue
		}
	}
}

// run is part of Start loop but was easily to write in separate function
func (svc *Service) run(ctx context.Context, sink sink) error {
	m := Metric{
		Tags: svc.options.Tags,
		Time: time.Now(),
	}
	w := nodeWtr(sink, m)
	var ignores []ignoreConstraint
	if len(svc.options.IgnorePaths) > 0 {
		var err error
		if ignores, err = compileIgnores(svc.options.IgnorePaths); err != nil {
			return err
		}
	}
	for modName := range svc.device.Modules() {
		if indexOf(svc.options.IgnoreModules, modName) >= 0 {
			continue
		}
		b, err := svc.device.Browser(modName)
		if err != nil {
			return err
		}
		root := b.RootWithContext(ctx)
		root, err = root.Constrain("content=nonconfig")
		for _, ignore := range ignores {
			root.Constraints.AddConstraint(ignore.name, 10, 10, ignore)
		}
		if err = root.UpdateInto(w); err != nil {
			return err
		}
	}
	return nil
}

// indexOf - find a string in list. There are tons of libs to pull in that
// solve this, but didn't want to add deps for just one function
func indexOf(candidates []string, target string) int {
	for i, c := range candidates {
		if c == target {
			return i
		}
	}
	return -1
}

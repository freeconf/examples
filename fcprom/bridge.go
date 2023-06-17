package fcprom

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"
	"time"

	"github.com/freeconf/restconf"
	"github.com/freeconf/restconf/device"
	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/val"
)

type Bridge struct {
	options       Options
	device        device.Device
	RenderMetrics RenderMetrics
	localServer   *http.Server
	Modules       Modules
}

type Modules struct {
	Ignore []string
}

func NewBridge(d device.Device) *Bridge {
	return &Bridge{
		device: d,
	}
}

type RenderMetrics struct {
	Duration time.Duration
	Count    int64
}

func (b *Bridge) Apply(options Options) error {
	if options.Port == "" {
		bwsr, err := b.device.Browser("fc-restconf")
		if err != nil {
			return err
		}
		if b == nil {
			return errors.New("no internal browser found and port not configured")
		}
		server, valid := bwsr.Root().Peek(nil).(*restconf.Server)
		if !valid {
			return errors.New("expected to find *restconf.Server when peeking at path 'restconf'")
		}
		var existingHandler = server.UnhandledRequestHandler
		server.UnhandledRequestHandler = func(w http.ResponseWriter, r *http.Request) {
			if r.RequestURI == "/metrics" {
				b.ServeHTTP(w, r)
			} else if existingHandler != nil {
				existingHandler(w, r)
			}
		}
	} else {
		if b.localServer != nil {
			b.localServer.Close()
		}
		demux := http.NewServeMux()
		demux.Handle("/metrics", b)
		b.localServer = &http.Server{
			Handler: demux,
			Addr:    options.Port,
		}
		go func() {
			if err := b.localServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("trouble with server server %s", err)
			}
		}()
	}
	b.options = options
	return nil
}

type Options struct {
	Port           string // ":2112"
	UseLocalServer bool
}

func (b *Bridge) Options() Options {
	return b.options
}

func (b *Bridge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.generate(w)
}

func (b *Bridge) generate(out io.Writer) error {
	e := newExporter()
	t0 := time.Now()

Modules:
	for name := range b.device.Modules() {
		for _, ignore := range b.Modules.Ignore {
			if ignore == name {
				continue Modules
			}
		}
		bwsr, err := b.device.Browser(name)
		if err != nil {
			return err
		}
		n := e.node(metricName(name), []string{})
		root := bwsr.Root()
		defer root.Release()
		sel, err := root.Constrain("content=nonconfig")
		if err != nil {
			return err
		}
		defer sel.Release()
		if err = sel.InsertInto(n); err != nil {
			return err
		}
		if err := writeMetrics(out, e.metrics); err != nil {
			return err
		}
	}
	b.RenderMetrics = RenderMetrics{
		Duration: time.Since(t0),
		Count:    int64(len(e.metrics)),
	}
	return nil
}

func writeMetrics(out io.Writer, metrics map[string]*metric) error {
	keys := make([]string, 0, len(metrics))
	for key := range metrics {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, id := range keys {
		m := metrics[id]
		if m.helpString != "" {
			io.WriteString(out, fmt.Sprintf("# HELP %s %s\n", id, m.helpString))
		}
		io.WriteString(out, fmt.Sprintf("# TYPE %s %s\n", id, m.metricType))
		if mv, isMultivariate := m.value.([]*multivariateValue); isMultivariate {
			for _, v := range mv {
				io.WriteString(out, id)
				labels := v.labels
				if len(m.labels) > 0 {
					labels = append(labels, m.labels...)
				}
				writeLabels(out, labels)
				io.WriteString(out, fmt.Sprintf(" %v\n", v.value))
			}
		} else {
			io.WriteString(out, id)
			writeLabels(out, m.labels)
			io.WriteString(out, fmt.Sprintf(" %v\n", m.value))
		}
	}
	return nil
}

func writeLabels(out io.Writer, labels []string) {
	if len(labels) > 0 {
		io.WriteString(out, " {")
		for i, label := range labels {
			if i > 0 {
				io.WriteString(out, ",")
			}
			io.WriteString(out, label)
		}
		io.WriteString(out, "}")
	}
}

var invalidChars = regexp.MustCompile("[-]")
var extraWhitespace = regexp.MustCompile(`\s+`)

func metricName(ident string) string {
	return invalidChars.ReplaceAllString(ident, "_")
}

func docString(desc string) string {
	return extraWhitespace.ReplaceAllString(desc, " ")
}

func convValue(v val.Value, f val.Format) interface{} {
	// TODO: support this
	if f.IsList() {
		return nil
	}

	switch f {
	case val.FmtBool:
		if v.(val.Bool) {
			return 1
		}
		return 0
	case val.FmtString:
		return nil
	}
	return v.Value()
}

type multivariateValue struct {
	labels []string
	value  interface{}
}

type metric struct {
	helpString string
	metricType string
	labels     []string
	value      interface{}
}

type exporter struct {
	metrics map[string]*metric
}

func newExporter() *exporter {
	return &exporter{
		metrics: make(map[string]*metric),
	}

}

// remove, meta. package added this
func findExtension(name string, candidates []*meta.Extension) *meta.Extension {
	for _, e := range candidates {
		if e.Ident() == name {
			return e
		}
	}
	return nil
}

func (e *exporter) add(id string, m meta.Meta, mvLabels []string, value interface{}) {
	thisMetric, found := e.metrics[id]
	if !found {
		metricType := "gauge"
		if ext := meta.FindExtension("counter", m.Extensions()); ext != nil {
			metricType = "counter"
		}
		thisMetric = &metric{
			metricType: metricType,
			helpString: docString(m.(meta.Describable).Description()),
		}
		e.metrics[id] = thisMetric
	}
	if len(mvLabels) > 0 {
		mv := multivariateValue{
			labels: mvLabels,
			value:  value,
		}
		if thisMetric.value == nil {
			thisMetric.value = []*multivariateValue{&mv}
		} else {
			thisMetric.value = append(thisMetric.value.([]*multivariateValue), &mv)
		}
	} else {
		thisMetric.value = value
	}
}

func (e *exporter) node(prefix string, mvLabels []string) node.Node {
	return &nodeutil.Basic{
		OnField: func(r node.FieldRequest, hnd *node.ValueHandle) error {
			value := convValue(hnd.Val, r.Meta.Type().Format())
			if value == nil {
				return nil
			}
			id := fmt.Sprintf("%s_%s", prefix, metricName(r.Meta.Ident()))
			e.add(id, r.Meta, mvLabels, value)
			return nil
		},
		OnChild: func(r node.ChildRequest) (node.Node, error) {
			if !r.New {
				return nil, nil
			}
			id := fmt.Sprintf("%s_%s", prefix, metricName(r.Meta.Ident()))
			if meta.IsList(r.Meta) {
				if findExtension("multivariate", r.Meta.Extensions()) != nil {
					if len(r.Meta.(*meta.List).KeyMeta()) == 0 {
						return nil, fmt.Errorf("multivariate definition for %s must have a key defined", r.Meta.Ident())
					}
					return e.multivariate(id, mvLabels), nil
				}
			}
			return e.node(id, mvLabels), nil
		},
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			if !r.New {
				return nil, nil, nil
			}
			id := fmt.Sprintf("%s_%d", prefix, r.Row)
			return e.node(id, mvLabels), nil, nil
		},
	}
}

func (e *exporter) multivariate(prefix string, mvLabels []string) node.Node {
	return &nodeutil.Basic{
		OnNext: func(r node.ListRequest) (node.Node, []val.Value, error) {
			if !r.New {
				return nil, nil, nil
			}
			listLabels := mvLabels
			metas := r.Meta.KeyMeta()
			for i, key := range r.Key {
				label := fmt.Sprintf("%s=\"%s\"", metas[i].Ident(), key.String())
				listLabels = append(listLabels, label)
			}
			return e.node(prefix, listLabels), r.Key, nil
		},
	}
}

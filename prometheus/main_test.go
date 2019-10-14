package basics_test

import (
	"fmt"
	"strings"

	"github.com/freeconf/bridge/prombridge"
	"github.com/freeconf/manage/device"
	"github.com/freeconf/manage/restconf"
	"github.com/freeconf/yang/source"

	"github.com/freeconf/examples/basics"
)

func Example_main() {

	// Normal setup found in ../basics
	car := basics.New()
	yangPath := source.Path(".:../yang")
	d := device.New(yangPath)
	if err := d.Add("car", basics.Manage(car)); err != nil {
		panic(err)
	}
	restconf.NewServer(d)

	// Here setup Prometheus bridge that will "scan" all local modules
	// and automatically export metrics to prometheus.
	bridge := prombridge.NewBridge(d)
	if err := d.Add("prom-bridge", prombridge.Manage(bridge)); err != nil {
		panic(err)
	}

	// We add the prom-bridge section to setup default integration with prometheus.
	// remember, Prometheus is pull model, so we don't need to know the prometheus
	// server(s)
	staticStartupConfig := `
		{
			"fc-restconf" : {
				"web" : {
					"port" : ":8080"
				}
			},
			"car" : {
				"speed" : 100
			},
			"prom-bridge" : {
				"service": {
					"useLocalServer" : true
				},
				"modules" : {
					"ignore" : ["ietf-yang-library"]
				}
			}
		}`
	if err := d.ApplyStartupConfig(strings.NewReader(staticStartupConfig)); err != nil {
		panic(err)
	}

	// trick to wait for cntrl-c... but you might have other services to start
	//   select {}

	// Output:
	// done
	fmt.Println("done")
}

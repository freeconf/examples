package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/freeconf/examples/car"
	"github.com/freeconf/examples/fcinflux"
	"github.com/freeconf/restconf"
	"github.com/freeconf/restconf/device"
	"github.com/freeconf/yang/source"
)

// Connect everything together into a server to start up
func main() {
	flag.Parse()

	// Your app here
	app := car.New()

	// where the yang files are stored
	ypath := source.Path("../../yang:../../car:../")

	// Device is just a container for browsers.  Needs to know where YANG files are stored
	d := device.New(ypath)

	// Device can hold multiple modules, here we are only adding one
	if err := d.Add("car", car.Manage(app)); err != nil {
		panic(err)
	}

	// Prometheus will look at all local modules which will be car and fc-restconf
	// unless configured to ignore the module
	s := fcinflux.NewService(d)
	if err := d.Add("fc-influx", fcinflux.Manage(s)); err != nil {
		panic(err)
	}

	// Select wire-protocol RESTCONF to serve the device.
	restconf.NewServer(d)

	// apply start-up config normally stored in a config file on disk
	config := fmt.Sprintf(`{
		"fc-restconf":{
			"debug": true,
			"web":{
				"port":":8090"
			}
		},
		"fc-influx" : {
			"options" : {
				"connection" : {
					"addr" : "http://localhost:8086",
					"apiToken" : "%s"
				},
				"organization": "freeconf",
				"database" : "demo",
				"bucket": "demo",
				"frequency": 10,
				"tags" : {
					"device": "abc123"
				},
				"ignoreModules": ["ietf-yang-library", "fc-restconf", "fc-influx"]
			}
		},
        "car":{"speed":10}
	}`, os.Getenv("API_TOKEN"))
	if err := d.ApplyStartupConfig(strings.NewReader(config)); err != nil {
		panic(err)
	}

	if !*testMode {
		// wait for ctrl-c
		log.Printf("server started")
		select {}
	}
}

var testMode = flag.Bool("test", false, "do not run in background (i.e. driven by unit test)")

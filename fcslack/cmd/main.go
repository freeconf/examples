package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/freeconf/examples/car"
	"github.com/freeconf/restconf"

	"github.com/freeconf/restconf/device"

	"github.com/freeconf/yang/source"

	"github.com/freeconf/examples/fcslack"
)

func main() {
	flag.Parse()

	// Your app here
	app := car.New()

	// where the yang files are stored
	ypath := source.Path("../../yang:../../car:..")

	// Device is just a container for browsers.  Needs to know where YANG files are stored
	d := device.New(ypath)

	// Device can hold multiple modules, here we are only adding one
	if err := d.Add("car", car.Manage(app)); err != nil {
		panic(err)
	}

	s := fcslack.NewService(d)

	if err := d.Add("fc-slack", fcslack.Manage(s)); err != nil {
		panic(err)
	}

	restconf.NewServer(d)

	// apply start-up config normally stored in a config file on disk
	config := fmt.Sprintf(`{
		"fc-restconf":{
			"debug": true,
			"web":{
				"port":":8090"
			}
		},
		"fc-slack" : {
			"client" : {
				"apiToken": "%s"
			},
			"subscription": [
				{
					"channel" : "%s",
					"module": "car",
					"path": "update"
				}
			]
		},
        "car":{"speed":100}
	}`, os.Getenv("SLACK_API_TOKEN"), os.Getenv("SLACK_CHANNEL"))

	// bootstrap config for all local modules
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

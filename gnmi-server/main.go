package main

import (
	"flag"
	"log"

	"github.com/freeconf/examples/car"
	"github.com/freeconf/gnmi"
	"github.com/freeconf/restconf/device"
	"github.com/freeconf/yang/source"
)

// Connect everything together into a server to start up
func main() {
	flag.Parse()

	// Your app here
	app := car.New()

	// where the yang files are stored
	ypath := source.Any(gnmi.InternalYPath, source.Path("../yang"))

	// Device is just a container for browsers.  Needs to know where YANG files are stored
	d := device.New(ypath)

	// Device can hold multiple modules, here we are only adding one
	if err := d.Add("car", car.Manage(app)); err != nil {
		panic(err)
	}

	// Select wire-protocol gNMI to serve the device
	gnmi.NewServer(d)

	// apply start-up config normally stored in a config file on disk
	if err := d.ApplyStartupConfigFile("./startup.json"); err != nil {
		panic(err)
	}

	if !*testMode {
		// wait for ctrl-c
		log.Printf("server started")
		select {}
	}
}

var testMode = flag.Bool("test", false, "do not run in background (i.e. driven by unit test)")

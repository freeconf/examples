package main

import (
	"flag"

	"github.com/freeconf/examples/garage"

	"github.com/freeconf/c2g/c2"
	"github.com/freeconf/c2g/device"
	"github.com/freeconf/c2g/meta"
	"github.com/freeconf/c2g/restconf"
)

var startup = flag.String("startup", "startup.json", "startup configuration file.")
var verbose = flag.Bool("verbose", false, "verbose")

func main() {
	flag.Parse()
	c2.DebugLog(*verbose)

	app := garage.NewGarage()

	// notice the garage doesn't need yang for car.  it will get
	// that from proxy, that will in turn get it from car node, having
	// said that, if it does find yang locally, it will use it
	yangPath := meta.PathStreamSource("..:../../../c2g/yang")

	d := device.New(yangPath)

	mgmt := restconf.NewServer(d)

	chkErr(d.Add("garage", garage.Manage(app)))

	// apply start-up config, just enough to initialize connection to
	// services that will finishing configuration
	chkErr(d.ApplyStartupConfigFile(*startup))

	var sub c2.Subscription
	if mgmt.CallHome != nil {
		mgmt.CallHome.OnRegister(func(d device.Device, u device.RegisterUpdate) {
			if sub != nil {
				sub.Close()
			}
			if u == device.Register {
				baseAddress := mgmt.CallHome.Options().Address
				dm := device.NewMapClient(d, baseAddress, restconf.ProtocolHandler(yangPath))
				sub = garage.ManageCars(app, dm)
			}
		})
	}

	// wait for cntrl-c...
	select {}
}

func chkErr(err error) {
	if err != nil {
		panic(err)
	}
}

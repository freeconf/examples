package main

import (
	"log"

	"github.com/freeconf/restconf"
	"github.com/freeconf/restconf/device"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/source"
)

// This is not specific to restconf, just example data structure for your
// application
type MyApp struct {
	Message string
}

func main() {
	// your app instance
	app := MyApp{}

	// where to find YANG files
	ypath := source.Any(source.Path("."), restconf.InternalYPath)

	// organize modules.
	d := device.New(ypath)

	// register your application module. you can register as many as you want here.
	//   param 1 - name of the module, "hello.yang" must exist in ypath
	//   param 2 - code that connects (bridges) from your App to yang interface
	//             there are many options in nodeutil package to base your
	//             implementation on.  Here we use reflection because our yang file aligns
	//             with out application data structure.
	rootNode := nodeutil.Reflect{}.Object(&app)
	d.Add("hello", rootNode)

	// select RESTCONF as management protocol. gNMI is option as well
	restconf.NewServer(d)

	// this will apply configuration and starting RESTCONF web server
	if err := d.ApplyStartupConfigFile("./startup.json"); err != nil {
		log.Fatal(err)
	}

	select {}
}

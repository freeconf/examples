package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/freeconf/examples/car"
	"github.com/freeconf/restconf"
	"github.com/freeconf/restconf/device"
	"github.com/freeconf/yang/source"
)

// Be sure to gen certs first. then use -secure=true
var secure = flag.Bool("secure", false, "startup with certs")

// Connect everything together into a server to start up
func main() {
	flag.Parse()

	// Your app here
	app := car.New()

	// where the yang files are stored
	ypath := source.Any(restconf.InternalYPath, source.Path("../yang:."))

	// Device is just a container for browsers.  Needs to know where YANG files are stored
	d := device.New(ypath)

	// Device can hold multiple modules, here we are only adding one
	if err := d.Add("car", car.Manage(app)); err != nil {
		panic(err)
	}

	// Example shows how to upload file. Remember, you can add as many modules as you want.
	if err := d.Add("file-upload", manageUploader()); err != nil {
		panic(err)
	}

	// Select wire-protocol RESTCONF to serve the device.
	s := restconf.NewServer(d)

	// This is the single line to inform RESTCONF server to serve your static
	// web site.
	s.RegisterWebApp(".", "index.html", "ui")

	// TLS not only is for security, but enables HTTP/2 which supports event stream
	// over primary http connection
	tls := ""
	if *secure {
		tls = `
			"tls": {
				"cert": {
					"certFile": "server.crt",
					"keyFile": "server.key"
				}
			},`
	}

	config := fmt.Sprintf(`{
		"fc-restconf":{
			"debug": true,
			"web":{
				%s
				"port":":8090",
				"writeTimeout": 0
			}
		},
		"car":{"speed":100}
	}`, tls)

	// apply start-up config normally stored in a config file on disk
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

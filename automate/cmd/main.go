package main

import (
	"flag"

	"github.com/freeconf/examples/automate"

	"github.com/freeconf/gconf/c2"
	"github.com/freeconf/gconf/device"
	"github.com/freeconf/gconf/meta"
)

// Initialize and start our RESTCONF proxy service.
//
// To run:
//    cd ./src/vendor/github.com/freeconf/gconf/examples/proxy
//    go run ./main.go
//
// Then open web browser to
//   http://localhost:8080/restconf/ui/index.html
//

var sysType = flag.String("sys", "go", "go is launches services this main.  Option : cmd - launches separate executables and a more complete test but much slower")
var verbose = flag.Bool("verbose", false, "logging on")

func main() {
	flag.Parse()
	c2.DebugLog(*verbose)
	var sys automate.System
	switch *sysType {
	case "go":
		sys = &automate.CmdSystem{
			ExamplesDir: " ../../",
			VarDir:      "./var",
			Verbose:     *verbose,
		}
	case "cmd":
		sys = &automate.GoSystem{
			Map:      device.NewMap(),
			YangPath: meta.PathStreamSource(".:../"),
		}
	default:
		panic(*sysType + " not a valid option")
	}

	automate.Garage(sys, 1)

	// Wait for cntrl-c...
	select {}
}

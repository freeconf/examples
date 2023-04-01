package demo

import (
	"fmt"

	"github.com/freeconf/restconf/client"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/source"
)

func connectClient() {

	// YANG: just need YANG file ietf-yang-library.yang, not the yang of remote system as that will
	// be downloaded as needed
	ypath := source.Path("../yang")

	// Connect
	proto := client.ProtocolHandler(ypath)
	dev, err := proto("http://localhost:9998/restconf")
	if err != nil {
		panic(err)
	}

	// Get a browser to walk server's management API for car
	car, err := dev.Browser("car")
	if err != nil {
		panic(err)
	}

	// Example of config: I feel the need, the need for speed
	// bad config is rejected in client before it is sent to server
	err = car.Root().UpsertFrom(nodeutil.ReadJSON(`{"speed":100}`)).LastErr
	if err != nil {
		panic(err)
	}

	// Example of metrics: Get all metrics as JSON
	metrics, err := nodeutil.WriteJSON(car.Root().Find("?content=nonconfig"))
	if err != nil {
		panic(err)
	}
	if metrics == "" {
		panic("no metrics")
	}

	// Example of RPC: Reset odometer
	err = car.Root().Find("reset").Action(nil).LastErr
	if err != nil {
		panic(err)
	}

	// Example of notification: Car has an important update
	unsub, err := car.Root().Find("update").Notifications(func(n node.Notification) {
		msg, err := nodeutil.WriteJSON(n.Event)
		if err != nil {
			panic(err)
		}
		fmt.Println(msg)
	})
	if err != nil {
		panic(err)
	}
	defer unsub()

	// Example of multiple modules: This is the FreeCONF server module
	rcServer, err := dev.Browser("fc-restconf")
	if err != nil {
		panic(err)
	}

	// Example of config: Enable debug logging on FreeCONF's remote RESTCONF server
	err = rcServer.Root().UpsertFrom(nodeutil.ReadJSON(`{"debug":true}`)).LastErr
	if err != nil {
		panic(err)
	}
}

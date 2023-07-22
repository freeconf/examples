#!/usr/bin/env python3

import freeconf.device
import freeconf.source
import freeconf.node
import freeconf.restconf
import freeconf.nodeutil.json
import threading

def connect_client():
	# YANG: just need YANG file ietf-yang-library.yang, not the yang of remote system as that will
	# be downloaded as needed
    ypath = freeconf.source.any(
        freeconf.source.restconf_internal_ypath(),
    )

    # connect
    dev = freeconf.device.Device.client(ypath, "http://localhost:9998/restconf")

    # Get a browser to walk server's management API for car
    car = dev.get_browser("car")
    root = car.root()

	# Example of config: I feel the need, the need for speed
	# bad config is rejected in client before it is sent to server
    root.upsert_from(freeconf.nodeutil.json.json_read_str('{"speed":100}'))

	# Example of metrics: Get all metrics as JSON
    sel = root.find("?content=nonconfig")
    metrics = freeconf.nodeutil.json.json_write_str(sel)
    if metrics == "":
        raise Exception("no metrics")
    print(metrics)

	# Example of RPC: Reset odometer
    reset = root.find("reset")
    reset.action(None)
    reset.release()

	# Example of notification: Car has an important update
    update_called = threading.Condition()
    update_called.acquire()
    def on_update(msg):
        nonlocal update_called
        update_called.notify()
    update = root.find("update")
    unsub = update.notification(on_update)
    print("waiting for update...")
    update_called.wait()
    unsub()

    root.release()

	# sel, err = root.Find("update")
	# if err != nil {
	# 	panic(err)
	# }
	# defer sel.Release()
	# unsub, err := sel.Notifications(func(n node.Notification) {
	# 	msg, err := nodeutil.WriteJSON(n.Event)
	# 	if err != nil {
	# 		panic(err)
	# 	}
	# 	fmt.Println(msg)
	# })
	# if err != nil {
	# 	panic(err)
	# }
	# defer unsub()

	# // Example of multiple modules: This is the FreeCONF server module
	# rcServer, err := dev.Browser("fc-restconf")
	# if err != nil {
	# 	panic(err)
	# }

	# // Example of config: Enable debug logging on FreeCONF's remote RESTCONF server
	# serverSel := rcServer.Root()
	# defer serverSel.Release()
	# serverSel.UpsertFrom(nodeutil.ReadJSON(`{"debug":true}`))
	# if err != nil {
	# 	panic(err)
	# }

import sys
sys.path.append('../python')
import car
import manage
import time

def start_server():
    ypath = freeconf.source.any(
        freeconf.source.restconf_internal_ypath(),
        freeconf.source.path("../yang")
    )
    mod = freeconf.parser.load_module_file(ypath, 'car')
    app = car.Car()
    mgmt = manage(app)
    b = freeconf.node.Browser(mod, mgmt)
    dev = freeconf.device.Device(ypath)
    dev.add_browser(b)

    _ = freeconf.restconf.Server(dev)
    dev.apply_startup_config_file("startup.json")
    print("server started")

start_server()
connect_client()
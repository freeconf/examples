#!/usr/bin/env python3
import sys
from freeconf import parser, device, source, node, gnmi
import threading

car_dir = "../car"

sys.path.append(car_dir)
import manage
import car

# where YANG files are
ypath = source.any(
    source.gnmi_internal_ypath(),
    source.path(car_dir)
)

mod = parser.load_module_file(ypath, 'car')

# create your empty app instance
app = car.Car()

# dev hosts al your modules
b = node.Browser(mod, manage(app))
dev = device.Device(ypath)
dev.add_browser(b)

_ = gnmi.Server(dev)
dev.apply_startup_config_file("startup.json")

# sleep forever
threading.Event().wait

from freeconf import restconf, source, device, parser, node, source, nodeutil
from threading import Event

# Represents a basic python application. There are no requirements imposed
# by FreeCONF on how you develop your application
class MyApp:
    def __init__(self):
        self.message = None

app = MyApp()

# The remaining is FreeCONF specific and shows how to build a management interface
# to your python application

# specify all the places where you store YANG files
ypath = source.any(
    source.path("."),                   # director to your local *.yang files
    source.restconf_internal_ypath()    # required for restconf protocol support
)

# load and validate your YANG file(s)
mod = parser.load_module_file(ypath, "hello")

# device hosts one or more management "modules" into a single instance that you
# want to export in the management interface
dev = device.Device(ypath)

# connect your application to your management implementation.
# there are endless ways to to build your management interface from code generation,
# to reflection and any combination there of.  A lot more information in docs.
mgmt = nodeutil.Node(app)

# connect parsed YANG to your management implementation.  Browser is a powerful way
# to dynamically control your application can can be useful in unit tests or other contexts
# but here we construct it to serve our management API
b = node.Browser(mod, mgmt)

# register your app management browser in device.  Device can hold any number of browsers
dev.add_browser(b)

# select RESTCONF as management protocol. gNMI is option as well or any custom or 
# future protocols
s = restconf.Server(dev)

# this will apply configuration including starting the RESTCONF web server
dev.apply_startup_config_file("./startup.json")

# simple python trick to wait until ctrl-c shutdown
Event().wait()
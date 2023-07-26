import freeconf.restconf
import freeconf.source
import freeconf.device
import freeconf.parser
import freeconf.nodeutil.reflect
from threading import Event

class MyApp:
    def __init__(self):
        self.message = None

app = MyApp()

# specify all the places where you store YANG files
ypath = freeconf.source.any(
    freeconf.source.path("."),
    freeconf.source.restconf_internal_ypath()
)

# load and validate your YANG file
mod = freeconf.parser.load_module_file(ypath, "hello")

# device hosts one or more management "modules" into a single instance that you
# want to export in management interface
dev = freeconf.device.Device(ypath)

# connect your application to your management implementation.
# there are endless ways to to build your management interface from code generation,
# to reflection and any combination there of.  A lot more information in docs.
mgmt = freeconf.nodeutil.reflect.Reflect(app)

# connect parsed YANG to your management implementation.
b = freeconf.node.Browser(mod, mgmt)

# register your  our app management 
dev.add_browser(b)

# select RESTCONF as management protocol. gNMI is option as well
s = freeconf.restconf.Server(dev)

# this will apply configuration including starting RESTCONF web server
dev.apply_startup_config("./startup.json")

# simple python trick to wait until shutdown
Event().wait()
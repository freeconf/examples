#!/usr/bin/env python3
import sys
import unittest 
import freeconf.driver
import freeconf.parser
import freeconf.node
import freeconf.device
import freeconf.restconf
import freeconf.source
import freeconf.nodeutil.json
import requests
import car
import manage

class TestManage(unittest.TestCase):

    def test_server(self):
        ypath = freeconf.source.any(
            freeconf.source.restconf_internal_ypath(),
            freeconf.source.path("../yang")
        )
        mod = freeconf.parser.load_module_file(ypath, 'car')
        app = car.Car()
        mgmt = manage.manage(app)
        b = freeconf.node.Browser(mod, mgmt)
        cfg = freeconf.nodeutil.json.json_write_str(b.root())
        print(cfg)


if __name__ == '__main__':
    unittest.main()

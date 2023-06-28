#!/usr/bin/env python3
import unittest 
import freeconf.parser
import freeconf.nodeutil
from manage import manage_bird
from bird import Bird 

class TestManage(unittest.TestCase):

    def test_manage(self):
        app = Bird("sparrow", 99, 1000)
        p = freeconf.parser.Parser()
        m = p.load_module('..', 'bird')
        mgmt = manage_bird(app)
        bwsr = freeconf.node.Browser(m, mgmt)
        root = bwsr.root()
        try:
            root.upsert_into(freeconf.nodeutil.json_write("tmp"))
        finally:
            root.release()
        with open("tmp", "r") as f:
            self.assertEqual('{"name":"sparrow","location":"99,1000"}', f.read())


if __name__ == '__main__':
    unittest.main()

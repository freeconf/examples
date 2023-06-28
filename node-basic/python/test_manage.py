#!/usr/bin/env python3
import unittest 
import freeconf.parser
import freeconf.nodeutil
from manage import ManageApp
from app import App 

class TestManage(unittest.TestCase):

    def test_manage(self):
        app = App()
        p = freeconf.parser.Parser()
        m = p.load_module('..', 'my-app')
        mgmt = ManageApp(app)
        bwsr = freeconf.node.Browser(m, mgmt)
        root = bwsr.root()
        try:
            root.upsert_into(freeconf.nodeutil.json_write("tmp"))
        finally:
            root.release()
        with open("tmp", "r") as f:
            self.assertEqual('{"users":{},"fonts":{},"bagels":{}}', f.read())


if __name__ == '__main__':
    unittest.main()

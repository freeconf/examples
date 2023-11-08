#!/usr/bin/env python3
import unittest 
from freeconf import parser, nodeutil, node, source
from manage import manage_app
from app import App 

class TestManage(unittest.TestCase):

    def test_manage(self):
        app = App()
        ypath = source.path("..")
        m = parser.load_module_file(ypath, 'my-app')
        mgmt = manage_app(app)
        bwsr = node.Browser(m, mgmt)
        root = bwsr.root()
        try:
            actual = nodeutil.json_write_str(root)
            self.assertEqual('{"users":{},"fonts":{},"bagels":{}}', actual)
        finally:
            root.release()

if __name__ == '__main__':
    unittest.main()

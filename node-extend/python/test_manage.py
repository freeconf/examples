#!/usr/bin/env python3
import unittest 
from freeconf import parser, nodeutil, source, node
from manage import manage_bird
from bird import Bird 

class TestManage(unittest.TestCase):

    def test_manage(self):
        app = Bird("sparrow", 99, 1000)
        ypath = source.path('..')
        m = parser.load_module_file(ypath, 'bird')
        mgmt = manage_bird(app)
        bwsr = node.Browser(m, mgmt)
        root = bwsr.root()
        try:
            actual = nodeutil.json_write_str(root)
            self.assertEqual('{"name":"sparrow","location":"99,1000"}', actual)
        finally:
            root.release()

if __name__ == '__main__':
    unittest.main()

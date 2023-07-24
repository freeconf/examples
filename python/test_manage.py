#!/usr/bin/env python3
import unittest 
from freeconf import parser, node, source, nodeutil
import manage
import queue
import car
import time

# Test the car management logic in manage.py
class TestManage(unittest.TestCase):

    def test_manage(self):
        ypath = source.any(
            source.restconf_internal_ypath(),
            source.path("../yang")
        )
        mod = parser.load_module_file(ypath, 'car')
        app = car.Car()

        # no web server needed, just your app and management function.
        brwsr = node.Browser(mod, manage.manage(app))
        root = brwsr.root()

        # read all the config
        curr_cfg = nodeutil.json_write_str(root.find("?content=config"))
        expected = '{"speed":1000,"pollInterval":1000,"tire":[{"pos":0,"size":"H15"},{"pos":1,"size":"H15"},{"pos":2,"size":"H15"},{"pos":3,"size":"H15"}]}'
        self.assertEqual(expected, curr_cfg)

        # access car and verify w/API
        self.assertEqual(False, app.running)

        # setup event listener, verify events later
        events = queue.Queue()
        def on_update(n):
            msg = nodeutil.json_write_str(n.event)
            events.put(msg)
        unsub = root.find("update").notification(on_update)

        # write config starts car
        root.update_from(nodeutil.json_read_str('{"speed":1000}'))
        self.assertEqual(1000, app.speed)

        # start car
        root.find("start").action()

        print(f"len after sub II {len(app.listeners)}")

        # verify first event 
        self.assertEqual('{"event":"carStarted"}', events.get())

        # done listening for events, ensure listener is removed
        unsub()

        #self.assertEqual(0, len(app.listeners))
        print(f"len after unsub {len(app.listeners)}")

        # hit all the RPCs
        root.find("rotateTires").action()
        root.find("replaceTires").action()
        root.find("reset").action()
        root.find("tire=0/replace").action()

        time.sleep(1)
        print(f"len after unsub II {len(app.listeners)}")


if __name__ == '__main__':
    unittest.main()

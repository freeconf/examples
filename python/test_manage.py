#!/usr/bin/env python3
import unittest 
from freeconf import parser, node, source, nodeutil
import manage
import queue
import car
import time

# Test the car management API in manage.py
class TestManage(unittest.TestCase):

    def test_manage(self):
        # where YANG files are
        ypath = source.path("../yang")
        mod = parser.load_module_file(ypath, 'car')

        # create your empty app instance
        app = car.Car()

        # no web server needed, just your app instance and management entry point.
        brwsr = node.Browser(mod, manage.manage(app))

        # get a selection into the management API to begin interaction
        root = brwsr.root()

        # TEST CONFIG GET: read all the config
        curr_cfg = nodeutil.json_write_str(root.find("?content=config"))
        expected = '{"speed":1000,"pollInterval":1000,"tire":[{"pos":0,"size":"H15"},{"pos":1,"size":"H15"},{"pos":2,"size":"H15"},{"pos":3,"size":"H15"}]}'
        self.assertEqual(expected, curr_cfg)

        # verify car starts not running.  Here we are checking from the car instance itsel
        self.assertEqual(False, app.running)

        # SETUP EVENT TEST: here we add a listener to receive events and stick them in a thread-safe
        # queue for assertion checking later
        events = queue.Queue()
        def on_update(n):
            msg = nodeutil.json_write_str(n.event)
            events.put(msg)
        unsub = root.find("update").notification(on_update)

        # TEST CONFIG SET: write config
        root.update_from(nodeutil.json_read_str('{"speed":2000}'))
        self.assertEqual(2000, app.speed)

        # TEST RPC: start car
        root.find("start").action()

        # TEST EVENT: verify first event. will block until event comes in
        print("waiting for event...")
        self.assertEqual('{"event":"carStarted"}', events.get())
        print("event receieved.")

        # TEST EVENT UNSUBSCRIBE: done listening for events
        unsub()

        # TEST RPCS: just hit all the RPCs for code coverage.  You could easily add the
        # underlying car object is changed accordingly
        root.find("rotateTires").action()
        root.find("replaceTires").action()
        root.find("reset").action()
        self.assertEqual(0, app.miles)
        root.find("tire=0/replace").action()
        print("done")

if __name__ == '__main__':
    unittest.main()

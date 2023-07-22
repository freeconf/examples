#!/usr/bin/env python3
import unittest 
import car
import queue

class TestCar(unittest.TestCase):

    def test_server(self):
        c = car.Car()
        c.poll_interval = 0.01
        events = queue.Queue()
        def on_update(e):
            events.put(e)
        unsub = c.on_update(on_update)
        print("waiting for car events...")
        c.start()
        self.assertEqual(car.EVENT_STARTED, events.get())
        self.assertEqual(car.EVENT_FLAT_TIRE, events.get())
        self.assertEqual(car.EVENT_STOPPED, events.get())
        c.replace_tires()
        self.assertEqual(car.EVENT_STARTED, events.get())
        unsub()
        c.stop()

        
if __name__ == '__main__':
    unittest.main()

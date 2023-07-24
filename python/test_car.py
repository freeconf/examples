#!/usr/bin/env python3
import unittest 
import queue
import car

class TestCar(unittest.TestCase):

    # Quick test of car's features using direct access to fields and methods
    def test_server(self):
        c = car.Car()
        c.poll_interval = 0.01
        c.speed = 1000

        events = queue.Queue()
        def on_update(e):
            print(f"got event {e}")
            events.put(e)
        unsub = c.on_update(on_update)
        print("waiting for car events...")
        c.start()
        
        self.assertEqual(car.EVENT_STARTED, events.get())
        self.assertEqual(car.EVENT_FLAT_TIRE, events.get())
        self.assertEqual(car.EVENT_STOPPED, events.get())
        c.replace_tires()
        c.start()

        self.assertEqual(car.EVENT_STARTED, events.get())
        unsub()

        c.replace_tires()

        c.stop()

        
if __name__ == '__main__':
    unittest.main()

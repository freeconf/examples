import freeconf.nodeutil
import time
import threading
import random

# If these strings match YANG enum ids then they will be converted automatically.  Ints would
# work as well
EVENT_STARTED = "carStarted"
EVENT_STOPPED = "carStopped"
EVENT_FLAT_TIRE = "flatTire"


# Simple application, no connection to management 
# C A R
# Your application code.
#
# Notice there are no reference to FreeCONF in this file.  This means your
# code remains:
# - unit test-able
# - Not auto-generated from model files
# - free of annotations/tags
class Car():

    def __init__(self):
        
        # metrics/state
        self.miles = 0
        self.running = False
        self.thread = None

        # potential configuable fields
        self.speed = 1000
        self.new_tires()
        self.poll_interval = 1.0 #secs

        # Listeners are common on manageable code.  Having said that, listeners
        # remain relevant to your application.  The manage.go file is responsible
        # for bridging the conversion from application to management api.
        self.listeners = []

    def start(self):
        if self.running:
            return
        self.thread = threading.Thread(target=self.run, name="Car")
        self.thread.start()

    def reset(self):
        self.miles = 0

    def stop(self):
        self.running = False

    def run(self):
        try:
            self.running = True
            self.update_listeners(EVENT_STARTED)
            while self.running:
                time.sleep(self.poll_interval)
                self.miles = self.miles + self.speed
                for t in self.tire:
                    t.endure_mileage(self.speed)
                    if t.flat:
                        self.update_listeners(EVENT_FLAT_TIRE)
                        return
        finally:
            self.running = False
            self.update_listeners(EVENT_STOPPED)

    def on_update(self, listener):
        self.listeners.append(listener)
        def closer():
            print("closed listener")
            self.listeners.remove(listener)
        return closer

    def update_listeners(self, event):
        for l in self.listeners:
            l(event)    

    def new_tires(self):
        self.tire = []
        for pos in range(4):
            self.tire.append(Tire(pos))
        self.last_rotation = self.miles

    def replace_tires(self):
        for t in self.tire:
            t.replace()
        self.last_rotation = int(self.miles)
        self.start()

    def rotate_tires(self):
        first = self.tire[0]
        self.tire[0] = self.tire[1]
        self.tire[1] = self.tire[2]
        self.tire[2] = self.tire[3]
        self.tire[3] = first
        for i in range(len(self.tire)):
            self.tire[i].pos = i
        self.last_rotation = int(self.miles)


class Tire:
    def __init__(self, pos):
        self.pos = pos
        self.wear = 100
        self.size = "H15"
        self.flat = False
        self.worn = False

    def replace(self):
        self.wear = 100
        self.flat = False
        self.worn = False

    def check_if_flat(self):
        if not self.flat:
            self.flat = (self.wear - (random.random() * 10)) < 0

    def endure_mileage(self, speed):
        # Wear down [0.0 - 0.5] of each tire proportionally to the tire position
        self.wear -= (float(speed) / 100) * float(self.pos) * random.random()
        self.check_if_flat()
        self.check_for_wear()    

    def check_for_wear(self):
        self.wear < 20

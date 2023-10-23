from freeconf import nodeutil, val

#
# C A R    M A N A G E M E N T
#  Bridge from model to car app.
#
# manage is root handler from car.yang. i.e. module car { ... }
def manage(c):

    # implement RPCs (action or rpc)
    def action(n, req):
        if req.meta.ident == 'reset':            
            c.miles = 0
            return None
        return n.do_action(req)
    
    # implement yang notifications which are really just events
    def notify(n, r):
        if r.meta.ident == 'update':
            def listener(event):
			    # events are nodes too
                r.send(nodeutil.Node(object={
                    "event": event
                }))
            closer = c.on_update(listener)
            return closer
        
        return n.do_notify(r)
    
    # implement fields that are not automatically handled by reflection.
    def read(_n, meta, v):
        if meta.units == 'millisecs':
            return int(v * 1000)
        return v

    def write(_n, meta, v):
        if meta.units == 'millisecs':
            return float(v) / 1000 # ms to secs
        return v

	# Extend and Reflect form a powerful combination, we're letting reflect do a lot of the CRUD
	# when the yang file matches the field names.  But we extend reflection
	# to add as much custom behavior as we want
    return nodeutil.Node(c,
        on_action = action, on_notify=notify, on_read=read, on_write=write)


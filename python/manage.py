from freeconf import nodeutil, val

#
# C A R    M A N A G E M E N T
#  Bridge from model to car app.
#
# manage is root handler from car.yang. i.e. module car { ... }
def manage(c):

    # implement navigation by containers and lists defined in yang file
    def child(p, req):
        if req.meta.ident == 'tire':
            return manage_tires(c)
        return p.child(req)

    # implement RPCs (action or rpc)
    def action(p, req):
        if req.meta.ident == 'stop':
            c.stop()
        elif req.meta.ident == 'start':
            c.start()
        elif req.meta.ident == 'rotateTires':            
            c.rotate_tires()
        elif req.meta.ident == 'replaceTires':            
            c.replace_tires()
        elif req.meta.ident == 'reset':            
            c.reset()
        else:
            return p.action(req)
        return None
    
    # implement yang notifications which are really just events
    def notification(p, r):
        if r.meta.ident == 'update':
            def listener(event):
			    # events are nodes too
                r.send(nodeutil.Reflect({
                    "event": event
                }))
            closer = c.on_update(listener)
            return closer
        
        return p.notification(r)
    
    # implement fields that are not automatically handled by reflection.
    def field(p, r, w):
        if r.meta.ident == 'pollInterval':
            if r.write:
                c.poll_interval = float(w.v) / 1000 # ms to secs
            else:
                return val.Val(int(c.poll_interval * 1000)) # secs to ms
        else:
            return p.field(r, w)
        return None

	# Extend and Reflect form a powerful combination, we're letting reflect do a lot of the CRUD
	# when the yang file matches the field names.  But we extend reflection
	# to add as much custom behavior as we want
    return nodeutil.Extend(
        base = nodeutil.Reflect(c),
        on_action = action, on_notification=notification, on_child=child, on_field=field)


# manage_tires handles list of tires.
def manage_tires(c):
    def next(r):
        key = r.key
        found = None
        if key:
            # request for a specific tire by key (pos)
            pos = key[0].v
            if pos < len(c.tire):
                found = c.tire[pos]
        elif r.row < len(c.tire):
            # request for the nth tire in list
            found = c.tire[r.row]
            key = [ val.Val(r.row) ]

        if found:
            return manage_tire(found), key
    return nodeutil.Basic(on_next=next)

# manage_tire handles each tire node.  Everything *inside* list tire { ... }
def manage_tire(t):
    def action(p, r):
        if r.meta.ident == "replace":
            t.replace()
        return None
	# again, let reflection do a lot of the work with one extension to handle replace tire
    # action
    return nodeutil.Extend(
        base=nodeutil.Reflect(t),
        on_action=action
    )

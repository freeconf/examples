import freeconf.nodeutil.reflect
import freeconf.nodeutil.extend
import freeconf.val

# Bridge car to FC management library
def manage(c):

    def action(node, req):
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
            return node.action(req)
        return None
    
    def child(node, req):
        if req.meta.ident == 'tire':
            return manage_tires(c)
        return node.child(req)

    def notification(node, req):
        if req.meta.ident == 'update':
            def listener(event):
                req.send(freeconf.nodeutil.reflect.Reflect({
                    "event": event
                }))
            closer = c.on_update(listener)
            return closer
        
        return node.notification(req)

    # because car's members and methods align with yang, we can use 
    # reflection for all of the CRUD
    return freeconf.nodeutil.extend.Extend(
        base = freeconf.nodeutil.reflect.Reflect(c),
        on_action = action, on_notification=notification, on_child=child)


def manage_tires(car):
    def next(r):
        key = r.key
        found = None
        if key:
            pos = key[0].value
            if pos < len(car.tire):
                found = car.tire[pos]
        elif r.row < len(car.tire):
            found = car.tire[r.row]
            key = [ freeconf.val.Val(freeconf.val.Format.INT32, r.row) ]

        if found:
            return manage_tire(found), key
    return freeconf.nodeutil.basic.Basic(on_next=next)


def manage_tire(t):
    def action(r):
        if r.meta.ident == "replace":
            t.replace()
        return None
    return freeconf.nodeutil.extend.Extend(
        base=freeconf.nodeutil.reflect.Reflect(t),
        on_action=action
    )

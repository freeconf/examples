import freeconf.nodeutil
import freeconf.val

def manage_bird(bird):
    base = freeconf.nodeutil.Reflect(bird)
    
    def on_field(parent, req, write_val):
        if req.meta.ident == "location":
            return freeconf.val.Val(freeconf.val.Format.STRING, bird.coordinates())
        return parent.field(req, write_val)
    
    return freeconf.nodeutil.Extend(base, on_field=on_field)
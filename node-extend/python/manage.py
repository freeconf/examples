from freeconf import nodeutil, val

def manage_bird(bird):
    base = nodeutil.Node(bird)
    
    def on_field(parent, req, write_val):
        if req.meta.ident == "location":
            return val.Val(bird.coordinates())
        return parent.field(req, write_val)
    
    return nodeutil.Extend(base, on_field=on_field)
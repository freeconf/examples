from freeconf import nodeutil

def manage_app(app):

    def child(req):
        if req.meta.ident == 'users':
            return nodeutil.Node(app.users)
        elif req.meta.ident == 'fonts':
            return nodeutil.Node(app.fonts)
        elif req.meta.ident == 'bagels':
            return nodeutil.Node(app.bagels)
        return None
    
    # while this could easily be nodeutil.Node, we illustrate a Basic
    # node should you want essentially an abstract class that stubs all 
    # this calls with reasonable default handlers
    return nodeutil.Basic(on_child=child)

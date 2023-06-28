import freeconf.nodeutil


class ManageApp(freeconf.nodeutil.Basic):

    def __init__(self, app):
        super().__init__()
        self.app = app        

    def child(self, req):
        if req.meta.ident == 'users':
            return freeconf.nodeutil.Reflect(self.app.users)
        elif req.meta.ident == 'fonts':
            return freeconf.nodeutil.Reflect(self.app.fonts)
        elif req.meta.ident == 'bagels':
            return freeconf.nodeutil.Reflect(self.app.bagels)
        return None

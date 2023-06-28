
class Bird():

    def __init__(self, name, x, y):
        self.name = name
        self.x = x
        self.y = y

    def coordinates(self):
        return f'{self.x},{self.y}'
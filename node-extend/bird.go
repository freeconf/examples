package demo

type Bird struct {
	Name     string
	Location Coordinates
}

type Coordinates struct{}

func (Coordinates) Set(string) {
}

func (Coordinates) Get() string {
	return "0,0"
}

package demo

import "fmt"

type Bird struct {
	Name string
	X    int
	Y    int
}

func (b *Bird) GetCoordinates() string {
	return fmt.Sprintf("%d,%d", b.X, b.Y)
}

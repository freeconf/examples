package chipmonk

type Friend struct {
	Name string
}

type Chipmunk struct {
	Friend map[string]*Friend
}

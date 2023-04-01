package demo

type App struct {
	users  *UserService
	fonts  *FontManager
	bagels *BagelMaker
}

type UserService struct{}
type FontManager struct{}
type BagelMaker struct{}

package screen

const ScreenWidth = 64
const ScreenHeight = 32

type Screen struct {
	screen [ScreenHeight][ScreenWidth]bool
}

func (s *Screen) Clear() {
	s.screen = [ScreenHeight][ScreenWidth]bool{}
}

func (s *Screen) Write(x uint8, y uint8, value bool) {
	s.screen[y][x] = value
}

func (s Screen) Read(x uint8, y uint8) bool {
	return s.screen[y][x]
}

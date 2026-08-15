package screen

import (
	"fmt"

	"github.com/Zyko0/go-sdl3/sdl"
)

const ScreenWidth = 64
const ScreenHeight = 32

const title = "Eit"
const scalingFactor = 20

type Screen struct {
	screen [ScreenHeight][ScreenWidth]bool

	window   *sdl.Window
	renderer *sdl.Renderer

	background sdl.Color
	pixel      sdl.Color
}

func New(width int, height int, background sdl.Color, pixel sdl.Color) (*Screen, error) {
	var s *Screen

	err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_GAMEPAD | sdl.INIT_AUDIO)

	if err != nil {
		return nil, fmt.Errorf("init SDL: %w", err)
	}

	window, renderer, err := sdl.CreateWindowAndRenderer(
		title,
		width*scalingFactor,
		height*scalingFactor,
		0,
	)

	if err != nil {
		return nil, fmt.Errorf("init window and renderer: %w", err)
	}

	err = renderer.SetLogicalPresentation(
		ScreenWidth,
		ScreenHeight,
		sdl.LOGICAL_PRESENTATION_INTEGER_SCALE,
	)

	if err != nil {
		return nil, fmt.Errorf("set renderer logical representation: %w", err)
	}

	s = &Screen{
		window:     window,
		renderer:   renderer,
		background: background,
		pixel:      pixel,
	}

	return s, nil
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

func (s *Screen) Update() {
	for y := range ScreenHeight {
		for x := range ScreenWidth {
			if s.Read(uint8(x), uint8(y)) {
				s.renderer.SetDrawColor(s.pixel.R, s.pixel.G, s.pixel.B, 0)
			} else {
				s.renderer.SetDrawColor(s.background.R, s.background.G, s.background.B, 0)
			}

			s.renderer.RenderPoint(float32(x), float32(y))
		}
	}

	s.renderer.Present()
}

func (s *Screen) Destroy() {
	sdl.Quit()

	s.window.Destroy()
	s.renderer.Destroy()
}

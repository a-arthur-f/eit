package keyboard

import "github.com/Zyko0/go-sdl3/sdl"

var KeyMapping = map[sdl.Keycode]uint8{
	sdl.K_1: 0x1,
	sdl.K_2: 0x2,
	sdl.K_3: 0x3,
	sdl.K_4: 0xc,
	sdl.K_Q: 0x4,
	sdl.K_W: 0x5,
	sdl.K_E: 0x6,
	sdl.K_R: 0xd,
	sdl.K_A: 0x7,
	sdl.K_S: 0x8,
	sdl.K_D: 0x9,
	sdl.K_F: 0xe,
	sdl.K_Z: 0xa,
	sdl.K_X: 0x0,
	sdl.K_C: 0xb,
	sdl.K_V: 0xf,
}

type Keyboard struct {
	keys [16]bool
}

func (k *Keyboard) Set(addr uint8, value bool) {
	k.keys[addr] = value
}

func (k Keyboard) Read(addr uint8) bool {
	return k.keys[addr]
}

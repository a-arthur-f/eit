package eit

import (
	"eit/cpu"
	"eit/font"
	"eit/keyboard"
	"eit/memory"
	"eit/screen"
	"fmt"
	"time"

	"github.com/Zyko0/go-sdl3/sdl"
)

const cpuFrequency = time.Second / 500
const timerFrequency = time.Second / 60

type Eit struct {
	mem *memory.Memory
	cpu *cpu.CPU
	scr *screen.Screen
	kbd *keyboard.Keyboard
}

func New() (Eit, error) {
	var eit Eit

	mem := memory.Memory{}
	kbd := keyboard.Keyboard{}

	loadFonts(&mem)

	background := sdl.Color{R: 18, G: 40, B: 76}
	pixel := sdl.Color{R: 138, G: 175, B: 234}

	scr, err := screen.New(64, 32, background, pixel)

	if err != nil {
		return eit, fmt.Errorf("init render: %w", err)
	}

	cpu := cpu.New(&mem, scr, &kbd)

	eit = Eit{
		cpu: &cpu,
		mem: &mem,
		scr: scr,
		kbd: &kbd,
	}

	return eit, nil
}

func (eit *Eit) Run() {
	lastCycle := time.Now()

	cpuAcc := time.Duration(0)
	timerAcc := time.Duration(0)

	running := true

	for running {
		var ev sdl.Event

		for sdl.PollEvent(&ev) {
			switch ev.Type {
			case sdl.EVENT_QUIT:
				running = false
			case sdl.EVENT_KEY_DOWN:
				key, ok := keyboard.KeyMapping[ev.KeyboardEvent().Key]

				if !ok {
					continue
				}
				
				eit.kbd.Set(key, true)
			case sdl.EVENT_KEY_UP:
				key, ok := keyboard.KeyMapping[ev.KeyboardEvent().Key]	
				
				if !ok {
					continue
				}

				eit.kbd.Set(key, false)

				if eit.cpu.WaitingInput() {
					eit.cpu.Input(key)
				}
			}
		}

		now := time.Now()
		elapsed := now.Sub(lastCycle)
		lastCycle = now
		
		cpuAcc += elapsed
		timerAcc += elapsed

		for cpuAcc >= cpuFrequency && !eit.cpu.WaitingInput() {
			eit.cpu.Cycle()
			cpuAcc -= cpuFrequency
		}

		for timerAcc >= timerFrequency {
			eit.cpu.TickTimers()
			timerAcc -= timerFrequency
		}

		if eit.cpu.ShouldDraw() {
			eit.scr.Update()
			eit.cpu.ResetDrawFlag()
		}
	}
}

func (eit Eit) Destroy() {
	eit.scr.Destroy()
}

func (eit *Eit) LoadRom(rom []byte) {
	address := uint16(0x200)

	for _, b := range rom {
		eit.mem.Write(address, b)
		address++
	}
}

func loadFonts(mem *memory.Memory) {
	fonts := font.Fonts()

	for ifont, font := range fonts {
		for ibyte, byte := range font  {
			mem.Write(uint16((ifont * len(font)) + ibyte), byte)
		}
	}
}

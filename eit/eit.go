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

const cpuFrequency = time.Second / 700
const timerFrequency = time.Second / 60

type Eit struct {
	mem    *memory.Memory
	cpu    *cpu.CPU
	scr    *screen.Screen
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
		cpu:    &cpu,
		mem:    &mem,
		scr:    scr,
	}

	return eit, nil
}

func (eit *Eit) Run() {
	lastCycle := time.Now()

	cpuAcc := time.Duration(0)
	timerAcc := time.Duration(0)

	running := true

	for running {
		now := time.Now()
		elapsed := now.Sub(lastCycle)
		lastCycle = now

		cpuAcc += elapsed
		timerAcc += elapsed

		for cpuAcc >= cpuFrequency {
			eit.cpu.Cycle()
			cpuAcc -= cpuFrequency
		}

		for timerAcc >= timerFrequency {
			eit.cpu.TickTimers()
			timerAcc -= timerFrequency
		}

		if eit.cpu.WaitingInput() {
			key := uint8(0x0)
			// handle input
			eit.cpu.Input(key)
		}

		var ev sdl.Event

		for sdl.PollEvent(&ev) {
			if ev.Type == sdl.EVENT_QUIT {
				running = false
			}
		}

		eit.scr.Update()
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
	for i, f := range font.Fonts() {
		for j, byte := range f {
			mem.Write(uint16(i*5+j), byte)
		}
	}
}

package eit

import (
	"eit/cpu"
	"eit/font"
	"eit/keyboard"
	"eit/memory"
	"eit/screen"
	"time"
)

const cpuFrequency = time.Second / 500
const timerFrequency = time.Second / 60

type Eit struct {
	mem *memory.Memory
	cpu *cpu.CPU
}

func New() Eit {
	mem := memory.Memory{}
	scr := screen.Screen{}
	kbd := keyboard.Keyboard{}

	loadFonts(&mem)

	cpu := cpu.New(&mem, &scr, &kbd)

	return Eit{
		mem: &mem,
		cpu: &cpu,
	}
}

func (eit *Eit) Run() {
	lastCycle := time.Now()

	cpuAcc := time.Duration(0)
	timerAcc := time.Duration(0)

	for {
		now := time.Now()
		elapsed := now.Sub(lastCycle)
		lastCycle = now

		cpuAcc += elapsed
		timerAcc += elapsed
		
		for cpuAcc >= cpuFrequency {
			eit.cpu.Cycle()	
			cpuAcc--
		}

		for timerAcc >= timerFrequency {
			eit.cpu.TickTimers()
			timerAcc--
		}

		if eit.cpu.WaitingInput() {
			key := uint8(0x0)
			// handle input
			eit.cpu.Input(key)
		}


		// handle graphics
	}
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
			mem.Write(uint16(i * 5 + j), byte)	
		}
	}
}

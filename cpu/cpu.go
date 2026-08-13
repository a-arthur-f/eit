package cpu

import (
	"fmt"
	"math/rand"

	"eit/font"
	"eit/keyboard"
	"eit/memory"
	"eit/screen"
)

type CPU struct {
	v  [16]uint8
	i  uint16
	pc uint16

	stack [16]uint16
	sp    uint8

	timer_delay uint8
	timer_sound uint8

	mem      *memory.Memory
	screen   *screen.Screen
	keyboard *keyboard.Keyboard

	waitingInput bool
}

func New(mem *memory.Memory, screen *screen.Screen, keyboard *keyboard.Keyboard) CPU {
	return CPU{
		pc:       0x200,
		mem:      mem,
		screen:   screen,
		keyboard: keyboard,
	}
}

func (cpu *CPU) Cycle() {
	if cpu.waitingInput {
		return
	}

	opcode := cpu.mem.Read(cpu.pc)
	cpu.execute(opcode)
}

func (cpu *CPU) Input(key uint8) {
	index := (cpu.getAddress() & 0xf00) >> 8
	cpu.v[index] = key
	
	cpu.waitingInput = false
}

func (cpu CPU) WaitingInput() bool {
	return cpu.waitingInput
}

func (cpu *CPU) execute(opcode uint8) {
	switch opcode & 0xf0 {
	case 0x00:
		instruction := cpu.mem.Read(cpu.pc + 1)

		switch instruction {
		case 0xe0:
			fmt.Println("clear")
			cpu.screen.Clear()
			cpu.pc += 2
		case 0xee:
			addr := cpu.pop()

			fmt.Printf("return to 0x%.4x\n", addr)
			cpu.pc = addr
		default:
			fmt.Printf("unknow instruction 0x%.2x for opcode 0x%.2x\n", instruction, opcode)
		}
	case 0x10:
		addr := cpu.getAddress()
		fmt.Printf("jmp to 0x%.4x\n", addr)
		cpu.pc = addr
	case 0x20:
		cpu.push(cpu.pc + 2)

		addr := cpu.getAddress()
		fmt.Printf("call 0x%.4x\n", addr)
		cpu.pc = addr
	case 0x30:
		index := cpu.mem.Read(cpu.pc) & 0x0f
		value := cpu.mem.Read(cpu.pc + 1)

		if cpu.v[index] != value {
			fmt.Println("no skip")
			cpu.pc += 2
		} else {
			fmt.Printf("skip to 0x%.4x\n", cpu.pc+4)
			cpu.pc += 4
		}
	case 0x40:
		index := cpu.mem.Read(cpu.pc) & 0x0f
		value := cpu.mem.Read(cpu.pc + 1)

		if cpu.v[index] == value {
			fmt.Println("no skip")
			cpu.pc += 2
		} else {
			fmt.Printf("skip to 0x%.4x\n", cpu.pc+4)
			cpu.pc += 4
		}
	case 0x50:
		x := cpu.mem.Read(cpu.pc) & 0x0f
		y := (cpu.mem.Read(cpu.pc+1) & 0xf0) >> 4

		if cpu.v[x] != cpu.v[y] {
			fmt.Println("no skip")
			cpu.pc += 2
		} else {
			fmt.Printf("skip to 0x%.4x\n", cpu.pc+4)
			cpu.pc += 4
		}
	case 0x60:
		addr := cpu.getAddress()

		index := addr & 0x0f00 >> 8
		value := addr & 0x00ff

		fmt.Printf("set V%x to %.2x\n", index, value)

		cpu.v[index] = uint8(value)

		cpu.pc += 2
	case 0x70:
		addr := cpu.getAddress()

		index := addr & 0x0f00 >> 8
		value := addr & 0x00ff

		fmt.Printf("add %.2x to V%x\n", value, index)

		cpu.v[index] += uint8(value)

		cpu.pc += 2
	case 0x80:
		instruction := cpu.mem.Read(cpu.pc+1) & 0x0f

		addr := cpu.getAddress()

		x := addr & 0x0f00 >> 8
		y := addr & 0x00f0 >> 4

		switch instruction {
		case 0x0:
			fmt.Printf("V[%d] = V[%d]\n", x, y)
			cpu.v[x] = cpu.v[y]
		case 0x1:
			fmt.Printf("V[%d] |= V[%d]\n", x, y)
			cpu.v[x] |= cpu.v[y]
		case 0x2:
			fmt.Printf("V[%d] &= V[%d]\n", x, y)
			cpu.v[x] &= cpu.v[y]
		case 0x3:
			fmt.Printf("V[%d] ^= V[%d]\n", x, y)
			cpu.v[x] ^= cpu.v[y]
		case 0x4:
			fmt.Printf("V[%d] += V[%d] with carry\n", x, y)

			carry := uint32(cpu.v[x])+uint32(cpu.v[y]) > 0xff

			if carry {
				cpu.v[0xf] = 0x1
			} else {
				cpu.v[0xf] = 0x0
			}

			cpu.v[x] += cpu.v[y]
		case 0x5:
			fmt.Printf("V[%d] -= V[%d] with borrow\n", x, y)

			borrow := cpu.v[x] < cpu.v[y]

			if borrow {
				cpu.v[0xf] = 0x0
			} else {
				cpu.v[0xf] = 0x1
			}

			cpu.v[x] -= cpu.v[y]
		case 0x6:
			fmt.Printf("V[%d] >>= V[%d]\n", x, y)

			leastBit := cpu.v[x] & 0x01
			cpu.v[0xf] = leastBit

			cpu.v[x] >>= cpu.v[y]
		case 0x7:
			fmt.Printf("V[%d] =- V[%d] with borrow\n", x, y)

			borrow := cpu.v[y] < cpu.v[x]

			if borrow {
				cpu.v[0xf] = 0x0
			} else {
				cpu.v[0xf] = 0x1
			}

			cpu.v[x] = cpu.v[y] - cpu.v[x]
		case 0xe:
			fmt.Printf("V[%d] <<= V[%d]\n", x, y)

			mostBit := cpu.v[x] & 0x80
			cpu.v[0xf] = mostBit >> 7

			cpu.v[x] <<= cpu.v[y]

		default:
			fmt.Printf("unknow instruction 0x%.1x for opcode 0x%.2x\n", instruction, opcode)
		}

		cpu.pc += 2
	case 0x90:
		x := cpu.mem.Read(cpu.pc) & 0x0f
		y := (cpu.mem.Read(cpu.pc+1) & 0xf0) >> 4

		if cpu.v[x] == cpu.v[y] {
			fmt.Println("no skip")
			cpu.pc += 2
		} else {
			fmt.Printf("skip to 0x%.4x\n", cpu.pc+4)
			cpu.pc += 4
		}
	case 0xa0:
		addr := cpu.getAddress() & 0xfff

		fmt.Printf("set I to 0x%.3x\n", addr)
		cpu.i = addr

		cpu.pc += 2
	case 0xb0:
		addr := (cpu.getAddress() & 0xfff) + uint16(cpu.v[0x0])

		fmt.Printf("jmp0 to 0x%.3x\n", addr)
		cpu.pc = addr
	case 0xc0:
		addr := cpu.getAddress()

		x := uint8((addr & 0xf00) >> 8)
		value := uint8(addr & 0x0ff)

		randNum := uint8(rand.Intn(0x100))

		fmt.Printf("set V[%d] to 0x%.2x\n", x, randNum&value)
		cpu.v[x] = randNum & value

		cpu.pc += 2
	case 0xd0:
		addr := cpu.getAddress()
		lines := addr & 0x00f
		xIndex := (addr & 0xf00) >> 8
		yIndex := (addr & 0x0f0) >> 4

		cpu.v[0xf] = 0x0

		for line := range lines {
			cols := cpu.mem.Read(cpu.i + line)

			for col := 0x7; col >= 0; col-- {
				x := cpu.v[xIndex] + (0x7-uint8(col))%screen.ScreenWidth
				y := (cpu.v[yIndex] + uint8(line)) % screen.ScreenHeight

				spriteBit := ((cols & uint8(1<<col)) >> col) != 0
				screenBit := cpu.screen.Read(x, y)

				fmt.Printf("set screen pixel(%d, %d) to %v\n", x, y, spriteBit != screenBit)
				cpu.screen.Write(x, y, spriteBit != screenBit)

				if spriteBit == screenBit {
					cpu.v[0xf] = 0x1
				}
			}
		}

		cpu.pc += 2
	case 0xe0:
		addr := cpu.getAddress()

		instruction := addr & 0x0ff
		index := (opcode & 0x0f) >> 4
		keyAddr := cpu.v[index]

		switch instruction {
		case 0x9e:
			if cpu.keyboard.Read(keyAddr) {
				cpu.pc += 4
			} else {
				cpu.pc += 2
			}
		case 0xa1:
			if cpu.keyboard.Read(keyAddr) {
				cpu.pc += 2
			} else {
				cpu.pc += 4
			}
		}
	case 0xf0:
		addr := cpu.getAddress()

		instruction := addr & 0x0ff
		index := addr & 0xf00 >> 8

		switch instruction {
		case 0x07:
			cpu.v[index] = cpu.timer_delay
		case 0x0a:
			cpu.waitingInput = true
		case 0x15:
			cpu.timer_delay = cpu.v[index]
		case 0x18:
			cpu.timer_sound = cpu.v[index]
		case 0x1e:
			cpu.i += uint16(cpu.v[index])
		case 0x29:
			cpu.i += font.Address(cpu.v[index])
		case 0x33:
			number := cpu.v[index]

			hundreds := number / 100
			tens := (number / 10) %10
			ones := number % 10

			cpu.mem.Write(cpu.i, hundreds)
			cpu.mem.Write(cpu.i + 1, tens)
			cpu.mem.Write(cpu.i + 2, ones)	
		}

		cpu.pc += 2
	default:
		fmt.Printf("unknow opcode 0x%.2x on adress 0x%.2x\n", opcode, cpu.pc)
	}
}

func (cpu CPU) getAddress() uint16 {
	return uint16(cpu.mem.Read(cpu.pc)&0x0f)<<8 | uint16(cpu.mem.Read(cpu.pc+1))
}

func (cpu *CPU) push(addr uint16) {
	cpu.stack[cpu.sp] = addr
	cpu.sp++
}

func (cpu *CPU) pop() uint16 {
	cpu.sp--
	return cpu.stack[cpu.sp]
}

func (cpu CPU) debug() {
	fmt.Printf("=============\n")
	fmt.Printf("v0: %.2x\n", cpu.v[0])
	fmt.Printf("v1: %.2x\n", cpu.v[1])
	fmt.Printf("v2: %.2x\n", cpu.v[2])
	fmt.Printf("v3: %.2x\n", cpu.v[3])
	fmt.Printf("v4: %.2x\n", cpu.v[4])
	fmt.Printf("v5: %.2x\n", cpu.v[5])
	fmt.Printf("v6: %.2x\n", cpu.v[6])
	fmt.Printf("v7: %.2x\n", cpu.v[7])
	fmt.Printf("v8: %.2x\n", cpu.v[8])
	fmt.Printf("v9: %.2x\n", cpu.v[9])
	fmt.Printf("va: %.2x\n", cpu.v[10])
	fmt.Printf("vb: %.2x\n", cpu.v[11])
	fmt.Printf("vc: %.2x\n", cpu.v[12])
	fmt.Printf("vd: %.2x\n", cpu.v[13])
	fmt.Printf("ve: %.2x\n", cpu.v[14])
	fmt.Printf("vf: %.2x\n", cpu.v[15])
	fmt.Println()
	fmt.Printf("stack: %v\n", cpu.stack)
	fmt.Printf("sp: %.4x\n", cpu.sp)
	fmt.Printf("=============\n")
}

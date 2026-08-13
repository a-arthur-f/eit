package cpu

import (
	"eit/font"
	"eit/keyboard"
	"eit/memory"
	"eit/screen"
	"reflect"
	"testing"
)

var scr = screen.Screen{}
var mem = memory.Memory{}
var key = keyboard.Keyboard{}

func Test00E0(t *testing.T) {
	opcode := uint16(0x00e0)
	loadOpcode(&mem, opcode)

	x := uint8(0x10)
	y := uint8(0x0a)

	scr.Write(x, y, true)

	cpu := New(&mem, &scr, &key)
	cpu.Cycle()

	expected := false
	got := scr.Read(x, y)

	if got != expected {
		t.Errorf("opcode 0x%.4x\ngot screen[0x10, 0x0a] %v\nwant %v", opcode, got, expected)
	}

	wantPc := uint16(0x202)

	if cpu.pc != wantPc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
	}
}

func Test00EE(t *testing.T) {
	opcode := uint16(0x00ee)
	loadOpcode(&mem, opcode)

	expectedPc := uint16(0x100)
	expectedSp := uint8(0)

	cpu := New(&mem, &scr, &key)
	cpu.push(expectedPc)

	cpu.Cycle()

	gotPc := cpu.pc
	gotSp := cpu.sp

	if gotPc != expectedPc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, gotPc, expectedPc)
	}

	if gotSp != expectedSp {
		t.Errorf("opcode 0x%.4x\ngot SP 0x%.4x\nwant 0x%.4x", opcode, gotSp, expectedSp)
	}
}

func Test1NNN(t *testing.T) {
	opcode := uint16(0x10ff)
	loadOpcode(&mem, opcode)

	expected := uint16(0x0ff)

	cpu := New(&mem, &scr, &key)
	cpu.Cycle()

	got := cpu.pc

	if got != expected {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, got, expected)
	}
}

func Test2NNN(t *testing.T) {
	opcode := uint16(0x20ff)
	loadOpcode(&mem, opcode)

	expectedPc := uint16(0x0ff)
	expectedSp := uint8(1)
	expectedStack := [16]uint16{0x202}

	cpu := New(&mem, &scr, &key)
	cpu.Cycle()

	gotPc := cpu.pc
	gotSp := cpu.sp
	gotStack := cpu.stack

	if gotPc != expectedPc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%0x4x\nwant 0x%.4x", opcode, gotPc, expectedPc)
	}

	if gotSp != expectedSp {
		t.Errorf("opcode 0x%.4x\ngot SP 0x%0x4x\nwant 0x%.4x", opcode, gotSp, expectedSp)
	}

	if !reflect.DeepEqual(gotStack, expectedStack) {
		t.Errorf("opcode 0x%.4x\ngot stack %v\nwant %v", opcode, gotStack, expectedStack)
	}
}

func Test3XNN(t *testing.T) {
	for x := range 15 {
		opcode := 0x3000 | uint16(x)<<8 | uint16(x)
		loadOpcode(&mem, opcode)

		cpu := New(&mem, &scr, &key)
		cpu.v[x] = uint8(x)
		cpu.Cycle()

		wantPC := 0x204

		if cpu.pc != uint16(wantPC) {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPC)
		}

		opcode += 2
		loadOpcode(&mem, opcode)

		cpu.pc = 0x200
		cpu.Cycle()

		wantPC = 0x202

		if cpu.pc != uint16(wantPC) {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPC)
		}
	}
}

func Test4XNN(t *testing.T) {
	for x := range 15 {
		opcode := 0x4000 | uint16(x)<<8 | uint16(x)
		loadOpcode(&mem, opcode)

		cpu := New(&mem, &scr, &key)
		cpu.v[x] = uint8(x)
		cpu.Cycle()

		wantPC := 0x202

		if cpu.pc != uint16(wantPC) {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPC)
		}

		opcode += 2
		loadOpcode(&mem, opcode)

		cpu.pc = 0x200
		cpu.Cycle()

		wantPC = 0x204

		if cpu.pc != uint16(wantPC) {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPC)
		}
	}
}

func Test5XY0(t *testing.T) {
	for x := range 15 {
		y := (x + 1) % 15
		opcode := 0x5000 | uint16(x)<<8 | uint16(y)<<4
		loadOpcode(&mem, opcode)

		cpu := New(&mem, &scr, &key)
		cpu.v[x] = uint8(x)
		cpu.v[y] = uint8(x)
		cpu.Cycle()

		wantPC := 0x204

		if cpu.pc != uint16(wantPC) {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPC)
		}

		loadOpcode(&mem, opcode)

		cpu.pc = 0x200
		cpu.v[y] = uint8(x + 2)
		cpu.Cycle()

		wantPC = 0x202

		if cpu.pc != uint16(wantPC) {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPC)
		}
	}
}

func Test6XNN(t *testing.T) {
	for x := range 15 {
		wantV := uint8(x * 2)

		opcode := 0x6000 | uint16(x)<<8 | uint16(wantV)
		loadOpcode(&mem, opcode)

		cpu := New(&mem, &scr, &key)
		cpu.Cycle()

		gotV := cpu.v[x]

		if gotV != wantV {
			t.Errorf("opcode 0x%.4x\ngot V[%d] 0x%.2x\nwant 0x%.2x", opcode, x, gotV, wantV)
		}

		wantPc := uint16(0x202)

		if cpu.pc != wantPc {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
		}
	}
}

func Test7XNN(t *testing.T) {
	for x := range 15 {
		opcode := 0x7000 | uint16(x)<<8 | uint16(x)
		loadOpcode(&mem, opcode)

		cpu := New(&mem, &scr, &key)
		cpu.v[x] = uint8(x * 2)

		wantV := cpu.v[x] + uint8(x)

		cpu.Cycle()

		gotV := cpu.v[x]

		if gotV != wantV {
			t.Errorf("opcode 0x%.4x\ngot V[%d] 0x%.2x\nwant 0x%.2x", opcode, x, gotV, wantV)
		}

		wantPc := uint16(0x202)

		if cpu.pc != wantPc {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
		}
	}
}

func Test8XY0(t *testing.T) {
	for x := range 15 {
		y := (x + 1) % 15

		opcode := 0x8000 | uint16(x)<<8 | uint16(y)<<4
		loadOpcode(&mem, opcode)

		cpu := New(&mem, &scr, &key)
		cpu.v[y] = uint8(y * 2)
		cpu.Cycle()

		if cpu.v[x] != cpu.v[y] {
			t.Errorf("opcode 0x%.4x\ngot V[%d] 0x%.2x\nwant to be == V[%d] == 0x%.2x", opcode, x, cpu.v[x], y, cpu.v[y])
		}

		wantPc := uint16(0x202)

		if cpu.pc != wantPc {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
		}
	}
}

func Test8XY1(t *testing.T) {
	for x := range 15 {
		y := (x + 1) % 15

		opcode := 0x8000 | uint16(x)<<8 | uint16(y)<<4 | 0x1
		loadOpcode(&mem, opcode)

		cpu := New(&mem, &scr, &key)
		cpu.v[x] = uint8(x * 3)
		cpu.v[y] = uint8(y)

		want := cpu.v[x] | cpu.v[y]

		cpu.Cycle()

		got := cpu.v[x]

		if got != want {
			t.Errorf("opcode 0x%.4x\ngot V[%d] 0x%.2x\nwant 0x%.2x", opcode, x, got, want)
		}

		wantPc := uint16(0x202)

		if cpu.pc != wantPc {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
		}
	}
}

func Test8XY2(t *testing.T) {
	for x := range 15 {
		y := (x + 1) % 15

		opcode := 0x8000 | uint16(x)<<8 | uint16(y)<<4 | 0x2
		loadOpcode(&mem, opcode)

		cpu := New(&mem, &scr, &key)
		cpu.v[x] = uint8(x * 3)
		cpu.v[y] = uint8(y)

		want := cpu.v[x] & cpu.v[y]

		cpu.Cycle()

		got := cpu.v[x]

		if got != want {
			t.Errorf("opcode 0x%.4x\ngot V[%d] 0x%.2x\nwant 0x%.2x", opcode, x, got, want)
		}

		wantPc := uint16(0x202)

		if cpu.pc != wantPc {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
		}
	}
}

func Test8XY3(t *testing.T) {
	for x := range 15 {
		y := (x + 1) % 15

		opcode := 0x8000 | uint16(x)<<8 | uint16(y)<<4 | 0x3
		loadOpcode(&mem, opcode)

		cpu := New(&mem, &scr, &key)
		cpu.v[x] = uint8(x * 3)
		cpu.v[y] = uint8(y)

		want := cpu.v[x] ^ cpu.v[y]

		cpu.Cycle()

		got := cpu.v[x]

		if got != want {
			t.Errorf("opcode 0x%.4x\ngot V[%d] 0x%.2x\nwant 0x%.2x", opcode, x, got, want)
		}

		wantPc := uint16(0x202)

		if cpu.pc != wantPc {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
		}
	}
}

func Test8XY4(t *testing.T) {
	for x := range 15 {
		y := (x + 1) % 15

		opcode := 0x8000 | uint16(x)<<8 | uint16(y)<<4 | 0x4
		loadOpcode(&mem, opcode)

		cpu := New(&mem, &scr, &key)
		cpu.v[x] = uint8(x)
		cpu.v[y] = uint8(y)

		want := cpu.v[x] + cpu.v[y]

		cpu.Cycle()

		got := cpu.v[x]

		if got != want {
			t.Errorf("opcode 0x%.4x\ngot V[%d] 0x%.2x\nwant 0x%.2x", opcode, x, got, want)
		}

		cpu.pc = 0x200
		cpu.v[x] = 0xff
		cpu.v[y] = 0xff

		cpu.Cycle()

		wantCarry := uint8(0x1)
		gotCarry := cpu.v[0xf]

		if gotCarry != wantCarry {
			t.Errorf("opcode 0x%.4x\ngot carry flag 0x%.1x\nwant 0x%.1x", opcode, gotCarry, wantCarry)
		}

		wantPc := uint16(0x202)

		if cpu.pc != wantPc {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
		}
	}
}

func Test8XY5(t *testing.T) {
	for x := range 15 {
		y := (x + 1) % 15

		opcode := 0x8000 | uint16(x)<<8 | uint16(y)<<4 | 0x5
		loadOpcode(&mem, opcode)

		cpu := New(&mem, &scr, &key)
		cpu.v[x] = uint8(x)
		cpu.v[y] = uint8(y)

		want := cpu.v[x] - cpu.v[y]

		cpu.Cycle()

		got := cpu.v[x]

		if got != want {
			t.Errorf("opcode 0x%.4x\ngot V[%d] 0x%.2x\nwant 0x%.2x", opcode, x, got, want)
		}

		cpu.pc = 0x200
		cpu.v[x] = 0x00
		cpu.v[y] = 0xff

		cpu.Cycle()

		wantBorrow := uint8(0x0)
		gotBorrow := cpu.v[0xf]

		if gotBorrow != wantBorrow {
			t.Errorf("opcode 0x%.4x\ngot borrow flag 0x%.1x\nwant 0x%.1x", opcode, gotBorrow, wantBorrow)
		}

		wantPc := uint16(0x202)

		if cpu.pc != wantPc {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
		}
	}
}

func Test8XY6(t *testing.T) {
	for x := range 15 {
		y := (x + 1) % 15

		opcode := 0x8000 | uint16(x)<<8 | uint16(y)<<4 | 0x6
		loadOpcode(&mem, opcode)

		cpu := New(&mem, &scr, &key)
		cpu.v[x] = uint8(x)
		cpu.v[y] = uint8(y)

		want := cpu.v[x] >> cpu.v[y]

		cpu.Cycle()

		got := cpu.v[x]

		if got != want {
			t.Errorf("opcode 0x%.4x\ngot V[%d] 0x%.2x\nwant 0x%.2x", opcode, x, got, want)
		}

		cpu.pc = 0x200
		cpu.v[x] = 0x01
		cpu.v[y] = 0xff

		cpu.Cycle()

		wantLeastBit := uint8(0x1)
		gotLeastBit := cpu.v[0xf]

		if gotLeastBit != wantLeastBit {
			t.Errorf("opcode 0x%.4x\ngot least significant bit 0x%.1x\nwant 0x%.1x", opcode, gotLeastBit, wantLeastBit)
		}

		cpu.pc = 0x200
		cpu.v[x] = 0x00
		cpu.v[y] = 0xff

		cpu.Cycle()

		wantLeastBit = uint8(0x0)
		gotLeastBit = cpu.v[0xf]

		if gotLeastBit != wantLeastBit {
			t.Errorf("opcode 0x%.4x\ngot least significant bit 0x%.1x\nwant 0x%.1x", opcode, gotLeastBit, wantLeastBit)
		}

		wantPc := uint16(0x202)

		if cpu.pc != wantPc {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
		}
	}
}

func Test8XY7(t *testing.T) {
	for x := range 15 {
		y := (x + 1) % 15

		opcode := 0x8000 | uint16(x)<<8 | uint16(y)<<4 | 0x7
		loadOpcode(&mem, opcode)

		cpu := New(&mem, &scr, &key)
		cpu.v[x] = uint8(x)
		cpu.v[y] = uint8(y)

		want := cpu.v[y] - cpu.v[x]

		cpu.Cycle()

		got := cpu.v[x]

		if got != want {
			t.Errorf("opcode 0x%.4x\ngot V[%d] 0x%.2x\nwant 0x%.2x", opcode, x, got, want)
		}

		cpu.pc = 0x200
		cpu.v[x] = 0xff
		cpu.v[y] = 0x00

		cpu.Cycle()

		wantBorrow := uint8(0x0)
		gotBorrow := cpu.v[0xf]

		if gotBorrow != wantBorrow {
			t.Errorf("opcode 0x%.4x\ngot borrow flag 0x%.1x\nwant 0x%.1x", opcode, gotBorrow, wantBorrow)
		}

		wantPc := uint16(0x202)

		if cpu.pc != wantPc {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
		}
	}
}

func Test8XYE(t *testing.T) {
	for x := range 15 {
		y := (x + 1) % 15

		opcode := 0x8000 | uint16(x)<<8 | uint16(y)<<4 | 0xe
		loadOpcode(&mem, opcode)

		cpu := New(&mem, &scr, &key)
		cpu.v[x] = uint8(x)
		cpu.v[y] = uint8(y)

		want := cpu.v[x] << cpu.v[y]

		cpu.Cycle()

		got := cpu.v[x]

		if got != want {
			t.Errorf("opcode 0x%.4x\ngot V[%d] 0x%.2x\nwant 0x%.2x", opcode, x, got, want)
		}

		cpu.pc = 0x200
		cpu.v[x] = 0xf0
		cpu.v[y] = 0xff

		cpu.Cycle()

		wantMostBit := uint8(0x1)
		gotMostBit := cpu.v[0xf]

		if gotMostBit != wantMostBit {
			t.Errorf("opcode 0x%.4x\ngot most significant bit 0x%.1x\nwant 0x%.1x", opcode, gotMostBit, wantMostBit)
		}

		cpu.pc = 0x200
		cpu.v[x] = 0x0f
		cpu.v[y] = 0xff

		cpu.Cycle()

		wantMostBit = uint8(0x0)
		gotMostBit = cpu.v[0xf]

		if gotMostBit != wantMostBit {
			t.Errorf("opcode 0x%.4x\ngot most significant bit 0x%.1x\nwant 0x%.1x", opcode, gotMostBit, wantMostBit)
		}

		wantPc := uint16(0x202)

		if cpu.pc != wantPc {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
		}
	}
}

func Test9XY0(t *testing.T) {
	for x := range 15 {
		y := (x + 1) % 15
		opcode := 0x9000 | uint16(x)<<8 | uint16(y)<<4
		loadOpcode(&mem, opcode)

		cpu := New(&mem, &scr, &key)
		cpu.v[x] = uint8(x)
		cpu.v[y] = uint8(x)
		cpu.Cycle()

		wantPC := 0x202

		if cpu.pc != uint16(wantPC) {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPC)
		}

		loadOpcode(&mem, opcode)

		cpu.pc = 0x200
		cpu.v[y] = uint8(x + 2)
		cpu.Cycle()

		wantPC = 0x204

		if cpu.pc != uint16(wantPC) {
			t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPC)
		}
	}
}

func TestANNN(t *testing.T) {
	want := uint16(0xf26)

	opcode := uint16(0xa000) | want
	loadOpcode(&mem, opcode)

	cpu := New(&mem, &scr, &key)
	cpu.Cycle()

	got := cpu.i

	if got != want {
		t.Errorf("opcode 0x%.4x\ngot I 0x%.4x\nwant 0x%.4x", opcode, got, want)
	}

	wantPc := uint16(0x202)

	if cpu.pc != wantPc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
	}
}

func TestBNNN(t *testing.T) {
	opcode := uint16(0xbf22)
	loadOpcode(&mem, opcode)

	cpu := New(&mem, &scr, &key)
	cpu.v[0x0] = 0x22

	wantPc := uint16(cpu.v[0x0]) + 0xf22

	cpu.Cycle()

	if wantPc != cpu.pc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
	}
}

func TestCXNN(t *testing.T) {
	opcode := uint16(0xc00f)
	loadOpcode(&mem, opcode)

	cpu := New(&mem, &scr, &key)
	cpu.Cycle()

	if cpu.v[0x0] > 0x0f {
		t.Errorf("opcode 0x%.4x\ngot value greater than 0x0f", opcode)
	}

	opcode = uint16(0xc0f0)
	loadOpcode(&mem, opcode)

	cpu.pc = 0x200
	cpu.Cycle()

	if cpu.v[0x0] > 0xf0 {
		t.Errorf("opcode 0x%.4x\ngot value greater than 0xf0", opcode)
	}

	opcode = uint16(0xc000)
	loadOpcode(&mem, opcode)

	cpu.pc = 0x200
	cpu.Cycle()

	if cpu.v[0x0] > 0x0 {
		t.Errorf("opcode 0x%.4x\ngot value greater than 0x00", opcode)
	}

	wantPc := uint16(0x202)

	if cpu.pc != wantPc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
	}

}

func TestDXYN(t *testing.T) {
	opcode := uint16(0xd018)
	loadOpcode(&mem, opcode)

	cpu := New(&mem, &scr, &key)

	cpu.v[0x0] = 27
	cpu.v[0x1] = 11

	cpu.mem.Write(0xe9f, 0b11111111)
	cpu.mem.Write(0xea0, 0b11111111)
	cpu.mem.Write(0xea1, 0b11111111)
	cpu.mem.Write(0xea2, 0b11111111)
	cpu.mem.Write(0xea3, 0b11111111)
	cpu.mem.Write(0xea4, 0b11111111)
	cpu.mem.Write(0xea5, 0b11111111)
	cpu.mem.Write(0xea6, 0b11111111)

	cpu.i = 0xe9f

	wantScreen := [screen.ScreenHeight][screen.ScreenWidth]bool{}

	for column := 11; column < 19; column++ {
		for row := 27; row < 35; row++ {
			wantScreen[column][row] = true
		}
	}

	cpu.screen.Clear()
	cpu.Cycle()

	for y := range uint8(32) {
		for x := range uint8(64) {
			screenPixel := cpu.screen.Read(x, y)
			wantPixel := wantScreen[y][x]

			if screenPixel != wantPixel {
				t.Errorf("wrong pixel(%d, %d)\ngot %v\nwant %v", x, y, screenPixel, wantPixel)
			}
		}
	}

	opcode = 0xd011
	loadOpcode(&mem, opcode)

	cpu.v[0x0] = 15
	cpu.v[0x1] = 0
	cpu.pc = 0x200

	cpu.Cycle()

	wantScreen[0][15] = true

	screenPixel := cpu.screen.Read(15, 0)
	wantPixel := wantScreen[0][15]

	if screenPixel != wantPixel {
		t.Errorf("wrong pixel(%d, %d)\ngot %v\nwant %v", 0, 15, screenPixel, wantPixel)
	}

	cpu.pc = 0x200
	cpu.Cycle()

	if cpu.v[0xf] != 0x1 {
		t.Errorf("opcode 0x%.4x\ngot V[0xF] 0x%.2x\nwant 0x1", opcode, cpu.v[0xf])
	}

	wantPc := uint16(0x202)

	if cpu.pc != wantPc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
	}
}

func TestEX9E(t *testing.T) {
	opcode := uint16(0xe09e)
	loadOpcode(&mem, opcode)

	cpu := New(&mem, &scr, &key)
	cpu.v[0] = 0x0
	cpu.Cycle()

	wantPc := uint16(0x202)

	if cpu.pc != wantPc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
	}

	cpu.keyboard.Set(0x0)
	cpu.pc = 0x200
	cpu.Cycle()

	wantPc = uint16(0x204)

	if cpu.pc != wantPc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
	}
}

func TestEXA1(t *testing.T) {
	opcode := uint16(0xe0a1)
	loadOpcode(&mem, opcode)

	cpu := New(&mem, &scr, &key)
	cpu.v[0] = 0x1
	cpu.Cycle()

	wantPc := uint16(0x204)

	if cpu.pc != wantPc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
	}

	cpu.keyboard.Set(0x1)
	cpu.pc = 0x200
	cpu.Cycle()

	wantPc = uint16(0x202)

	if cpu.pc != wantPc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
	}
}

func TestFX07(t *testing.T) {
	opcode := uint16(0xf207)
	loadOpcode(&mem, opcode)

	cpu := New(&mem, &scr, &key)
	cpu.timer_delay = 0x8f
	cpu.Cycle()

	want := cpu.timer_delay
	got := cpu.v[0x2]

	if got != want {
		t.Errorf("opcode 0x%.4x\ngot V[%d] 0x%.2x\nwant 0x%.2x", opcode, 2, got, want)
	}

	wantPc := uint16(0x202)

	if cpu.pc != wantPc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
	}
}

func TestFX0A(t *testing.T) {
	opcode := uint16(0xf00a)
	loadOpcode(&mem, opcode)

	cpu := New(&mem, &scr, &key)
	cpu.Cycle()
	cpu.Cycle()
	cpu.Cycle()

	wantPc := uint16(0x202)

	if cpu.pc != wantPc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
	}

	wantWaiting := true
	gotWaiting := cpu.waitingInput

	if gotWaiting != wantWaiting {
		t.Errorf("opcode 0x%.4x\ngot waiting == %v\nwant %v", opcode, gotWaiting, wantWaiting)
	}

	cpu.Input(0xf)

	wantWaiting = false
	gotWaiting = cpu.waitingInput

	if gotWaiting != wantWaiting {
		t.Errorf("opcode 0x%.4x\ngot waiting == %v\nwant %v", opcode, gotWaiting, wantWaiting)
	}

	want := uint8(0xf)
	got := cpu.v[0x0]

	if got != want {
		t.Errorf("opcode 0x%.4x\ngot V[%d] 0x%.2x\nwant 0x%.2x", opcode, 0x0, got, want)
	}
}

func TestFX15(t *testing.T) {
	opcode := uint16(0xf015)
	loadOpcode(&mem, opcode)

	cpu := New(&mem, &scr, &key)
	cpu.v[0x0] = 0xff
	cpu.Cycle()

	want := uint8(0xff)
	got := cpu.timer_delay

	if got != want {
		t.Errorf("opcode 0x%.4x\ngot Delay 0x%.2x\nwant 0x%.2x", opcode, got, want)
	}

	wantPc := uint16(0x202)

	if cpu.pc != wantPc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
	}
}

func TestFX18(t *testing.T) {
	opcode := uint16(0xf018)
	loadOpcode(&mem, opcode)

	cpu := New(&mem, &scr, &key)
	cpu.v[0x0] = 0xff
	cpu.Cycle()

	want := uint8(0xff)
	got := cpu.timer_sound

	if got != want {
		t.Errorf("opcode 0x%.4x\ngot Delay 0x%.2x\nwant 0x%.2x", opcode, got, want)
	}

	wantPc := uint16(0x202)

	if cpu.pc != wantPc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
	}
}

func TestFX1E(t *testing.T) {
	opcode := uint16(0xf01e)
	loadOpcode(&mem, opcode)

	cpu := New(&mem, &scr, &key)
	cpu.v[0x0] = 0xf
	cpu.i = 0xff
	cpu.Cycle()

	want := uint16(0xff + 0xf)
	got := cpu.i

	if got != want {
		t.Errorf("opcode 0x%.4x\ngot I 0x%.4x\nwant 0x%.4x", opcode, got, want)
	}

	wantPc := uint16(0x202)

	if cpu.pc != wantPc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
	}
}

func TestFX29(t *testing.T) {
	opcode := uint16(0xf029)
	loadOpcode(&mem, opcode)

	cpu := New(&mem, &scr, &key)
	cpu.v[0x0] = 0xf
	cpu.Cycle()

	addr := font.Address(0xf)

	want := addr
	got := cpu.i

	if got != want {
		t.Errorf("opcode 0x%.4x\ngot I 0x%.4x\nwant 0x%.4x", opcode, got, want)
	}

	wantPc := uint16(0x202)

	if cpu.pc != wantPc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
	}
}

func TestFX33(t *testing.T) {
	opcode := uint16(0xf033)
	loadOpcode(&mem, opcode)

	tests := []struct {
		name  string
		value uint8
		want1 uint8
		want2 uint8
		want3 uint8
	}{
		{name: "255", value: 0xff, want1: 2, want2: 5, want3: 5},
		{name: "15", value: 0xf, want1: 0, want2: 1, want3: 5},
		{name: "5", value: 0x5, want1: 0, want2: 0, want3: 5},
	}

	cpu := New(&mem, &scr, &key)
	cpu.i = 0x50

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu.pc = 0x200
			cpu.v[0x0] = tt.value
			cpu.Cycle()

			got1 := cpu.mem.Read(cpu.i)
			got2 := cpu.mem.Read(cpu.i + 1)
			got3 := cpu.mem.Read(cpu.i + 2)

			if got1 != tt.want1 {
				t.Errorf("opcode 0x%.4x\ngot first digit %d\nwant %d", opcode, got1, tt.want1)
			}

			if got2 != tt.want2 {
				t.Errorf("opcode 0x%.4x\ngot second digit %d\nwant %d", opcode, got3, tt.want3)
			}

			if got3 != tt.want3 {
				t.Errorf("opcode 0x%.4x\ngot third digit %d\nwant %d", opcode, got3, tt.want3)
			}
		})
	}

	wantPc := uint16(0x202)

	if cpu.pc != wantPc {
		t.Errorf("opcode 0x%.4x\ngot PC 0x%.4x\nwant 0x%.4x", opcode, cpu.pc, wantPc)
	}
}

func TestLoadOpcode(t *testing.T) {
	loadOpcode(&mem, 0x20ff)

	gotHighByte := mem.Read(0x200)
	gotLowByte := mem.Read(0x201)

	if gotHighByte != 0x20 {
		t.Errorf("mem[0x200]\ngot %.2x\nwant 0x20", gotHighByte)
	}

	if gotLowByte != 0xff {
		t.Errorf("mem[0x201]\ngot %.2x\nwant 0xff", gotLowByte)
	}
}

func loadOpcode(mem *memory.Memory, opcode uint16) {
	mem.Write(0x200, uint8(opcode>>8))
	mem.Write(0x201, uint8(opcode&0xff))
}

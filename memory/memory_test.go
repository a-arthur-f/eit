package memory

import (
	"testing"
)

func TestRead(t *testing.T) {
	em := Memory{}
	
	addr := uint16(0xff)
	expected := uint8(0x08)
	em.mem[addr] = expected

	got := em.Read(addr)

	if got != expected {
		t.Errorf("Read(0x%.4x)\ngot 0x%.2x\nwant 0x%.2x", addr, got, expected)
	}
}

func TestWrite(t *testing.T) {
	em := Memory{}

	addr := uint16(0x3f0)
	data := uint8(0x02)
	expected := uint8(0x02)

	em.Write(addr, data)

	got := em.mem[addr]

	if got != expected {
		t.Errorf("Write(0x%x4x, 0x%.2x)\ngot 0%x.2x\nwant 0x%.2x", addr, data, got, expected)	
	}
}

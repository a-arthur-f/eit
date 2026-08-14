package eit

import (
	"eit/font"
	"testing"

	"github.com/Zyko0/go-sdl3/bin/binsdl"
)

func TestInitialize(t *testing.T) {
	defer binsdl.Load().Unload()

	eit, err := New()

	if err != nil {
		t.Errorf("failed to initialize: %v", err)
	}

	fonts := font.Fonts()
	totalBytes := 16 * 5

	for i := 0; i < totalBytes; i += 5 {
		for j := range 5 {
			memByte := eit.mem.Read(uint16(i + j))
			fontByte := fonts[i / 5][j]

			if memByte != fontByte {
				t.Errorf("wrong font byte at mem[0x%.4x]\ngot 0x%.4x\nwant 0x%.4x", uint16(i + j), memByte, fontByte)
			}
		}
	}

	eit.Destroy()
}

func TestLoadRom(t *testing.T) {
	defer binsdl.Load().Unload()

	eit, err := New()

	if err != nil {
		t.Errorf("failed to initialize: %v", err)
	}

	rom := []byte{0x00, 0xff, 0xca, 0x44}
	eit.LoadRom(rom)

	for i, b := range rom {
		memByte := eit.mem.Read(uint16(0x200 + i))

		if memByte != b {
			t.Errorf("wrong rom byte loaded at mem[0x%.4x]\ngot 0x%.4x\nwant 0x%.4x", 0x200 + i, b, memByte)
		}
	}

	eit.Destroy()
}

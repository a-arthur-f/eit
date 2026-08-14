package main

import (
	"eit/eit"
	"log"
	"os"

	"github.com/Zyko0/go-sdl3/bin/binsdl"
)

func main() {
	defer binsdl.Load().Unload()

	eit, err := eit.New()

	if err != nil {
		log.Fatalf("Failed to init Eit: %v", err)
	}

	defer eit.Destroy()

	rom, err := os.ReadFile("./roms/1-chip8-logo.ch8")

	if err != nil {
		log.Fatalf("Failed to read rom: %v", err)
	}

	eit.LoadRom(rom)
	eit.Run()
}

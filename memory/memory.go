package memory

type Memory struct {
	mem [4096]uint8
}

func (mem Memory) Read(address uint16) uint8 {
	return mem.mem[address]
}

func (mem *Memory) Write(address uint16, data uint8) {
	mem.mem[address] = data
}

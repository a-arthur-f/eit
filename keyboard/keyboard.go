package keyboard

type Keyboard struct {
	keys [16]bool
}

func (k *Keyboard) Set(addr uint8) {
	k.keys[addr] = !k.keys[addr]
}

func (k Keyboard) Read(addr uint8) bool {
	return k.keys[addr]
}

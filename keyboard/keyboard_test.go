package keyboard

import (
	"testing"
)

func TestSet(t *testing.T) {
	key := Keyboard{}
	ki := uint8(0x8)

	want := true
	key.Set(ki, true)
	got := key.keys[ki]

	if got != want {
		t.Errorf("set key 0x%.1x\ngot %v\nwant %v", ki, got, want)
	}

	want = false
	key.Set(ki, false)
	got = key.keys[ki]

	if got != want {
		t.Errorf("set key 0x%.1x\ngot %v\nwant %v", ki, got, want)
	}
}

func TestRead(t *testing.T) {
	key := Keyboard{}
	ki := uint8(0xf)

	key.keys[ki] = true

	want := true
	got  := key.Read(ki)

	if got != want {
		t.Errorf("read key 0x%.1x\ngot %v\nwant %v", ki, got, want)
	}
}

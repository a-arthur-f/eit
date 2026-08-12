package font

import (
	"testing"
)

func TestAddress(t *testing.T) {
	tests := []struct {
		name  string
		value uint8
		want  uint16
	}{
		{name: "0x0", value: 0x0, want: 0},
		{name: "0x5", value: 0x5, want: 25},
		{name: "0xf", value: 0xf, want: 75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Address(tt.value)

			if got != tt.want {
				t.Errorf("Address(0x%.2x)\ngot 0x%.2x\nwant 0x%.2x", tt.value, got, tt.want)
			}
		})
	}
}

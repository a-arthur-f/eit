package screen

import (
	"reflect"
	"testing"
)

func TestClear(t *testing.T) {
	s := Screen{
		screen: [screenHeight][screenWidth]bool{
			{true, true, false, false}, 
			{true, true, true},
		},
	}

	expected := [screenHeight][screenWidth]bool{}

	s.Clear()

	if !reflect.DeepEqual(s.screen, expected) {
		t.Errorf("s.screen = %v\nwant %v", s.screen, expected)
	}
}

func TestWrite(t *testing.T) {
	s := Screen{}
	s.screen[10][40] = true

	tests := []struct {
		name string
		x, y uint8
		value bool
		expected bool
	}{
		{name: "true", x: 50, y: 20, value: true, expected: true},

		{name: "false", x: 40, y: 10, value: false, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.Write(tt.x, tt.y, tt.value)

			got := s.screen[tt.y][tt.x]

			if got != tt.expected {
				t.Errorf("Write(%d, %d, %v)\ngot %v\nwant %v", tt.x, tt.y, tt.value, got, tt.expected)
			}
		})
	}
}

func TestRead(t *testing.T) {
	s := Screen{
		screen: [screenHeight][screenWidth]bool{{false, true}},
	}

	tests := []struct {
		name string
		x, y uint8
		expected bool
	}{
		{name: "x: 1, y: 0", x: 1, y: 0, expected: true},
		{name: "x: 0, y: 0", x: 0, y: 0, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.Read(tt.x, tt.y)

			if got != tt.expected {
				t.Errorf("Read(%d, %d)\ngot %v\nwant %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}

package main

import "testing"

func TestAdd(t *testing.T) {
	got := add(1, 2)
	want := 3

	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{name: "2*3", a: 2, b: 3, want: 6},
		{name: "0*10", a: 0, b: 10, want: 0},
		{name: "-2*4", a: -2, b: 4, want: -8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := multiply(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

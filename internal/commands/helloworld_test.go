package commands

import "testing"

func TestGreet(t *testing.T) {
	got := Greet()
	want := "Hello, world!"

	if got != want {
		t.Errorf("Greet() = %q, want %q", got, want)
	}
}

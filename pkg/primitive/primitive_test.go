package primitive

import "testing"

func TestNew(t *testing.T) {
	workers := 1
	expected := Config{
		workers:    workers,
		OutputSize: 1280,
		Shape:      ShapeAny,
		Iterations: 200,
		Repeat:     1,
		Alpha:      128,
		Extension:  "jpg",
	}

	c := New(workers)
	if c != expected {
		t.Errorf("Got %+v;\n want: %+v", c, expected)
	}
}

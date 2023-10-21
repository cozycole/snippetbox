package assert

import (
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	// Called so that when t.Errorf() is called
	// from Equal() function, the Go test runner will
	// report the filename/line num of the code which called our Equal()
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v,; want: %v", actual, expected)
	}

}

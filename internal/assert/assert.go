package assert

import (
	"strings"
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

func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}

func NilError(t *testing.T, actual error) {
	t.Helper()

	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
	}
}

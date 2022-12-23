package util

import (
	"testing"
)

func TestAfter(t *testing.T) {
	fullString := "123 this is 123 description"
	key := "123"
	want := " this is 123 description"
	got := After(fullString, key)
	if want != got {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

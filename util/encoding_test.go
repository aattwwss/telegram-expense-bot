package util

import (
	"testing"
)

func TestToJson(t *testing.T) {
	type testStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	input := testStruct{Name: "Alice", Age: 30}
	got, err := ToJson(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `{"name":"Alice","age":30}`
	if got != want {
		t.Errorf("ToJson() = %q, want %q", got, want)
	}
}

func TestToJsonEmpty(t *testing.T) {
	type emptyStruct struct{}
	got, err := ToJson(emptyStruct{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "{}" {
		t.Errorf("ToJson() = %q, want {}", got)
	}
}

package main

import (
	"testing"
)

func TestUnpackValid(t *testing.T) {
	dataValid := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "Empty", input: "", expected: ""},
		{name: "Nothing to unpack", input: "Simple string", expected: "Simple string"},
		{name: "Simple unpack", input: "Unpack me dad2y plz3", expected: "Unpack me daddy plzzz"},
		{name: "Multirune character", input: "e\u03015", expected: "e\u0301\u0301\u0301\u0301\u0301"},
		{name: "Escaped #1", input: "qwe\\4\\5", expected: "qwe45"},
		{name: "Escaped #2", input: "qwe\\14", expected: "qwe1111"},
		{name: "Escaped #3", input: "qwe\\\\5", expected: "qwe\\\\\\\\\\"},
	}

	for _, data := range dataValid {
		t.Run(data.name, func(t *testing.T) {
			actual, err := unpack(data.input)

			if err != nil {
				t.Fatalf("Unexpected error: '%q'\n", err.Error())
			}

			if actual != data.expected {
				t.Fatalf("Actual result '%+q' is not equal to expected '%+q'\n", actual, data.expected)
			}
		})
	}
}

func TestUnpackInvalid(t *testing.T) {
	dataInvalid := []struct {
		name  string
		input string
		error string
	}{
		{name: "Only numbers", input: "12345", error: "String must begin with a letter"},
		{name: "Starts with escaped number", input: "\\12345", error: "String must begin with a letter"},
	}

	for _, data := range dataInvalid {
		t.Run(data.name, func(t *testing.T) {
			actual, err := unpack(data.input)

			if actual != "" {
				t.Fatalf("Unexpected returned string: '%s'\n", actual)
			}

			if err == nil {
				t.Fatalf("Expected error '%s', got nil\n", data.error)
			}

			if err.Error() != data.error {
				t.Fatalf("Actual error '%s' is not equal to expected '%s'\n", err.Error(), data.error)
			}
		})
	}
}

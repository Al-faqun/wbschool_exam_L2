package main

import (
	"testing"
)

// todo: add cases with escaping characters

func TestUnpackValid(t *testing.T) {
	dataValid := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "Empty", input: "", expected: ""},
		{name: "Nothing to unpack", input: "Simple string", expected: "Simple string"},
		{name: "Simple unpack", input: "Unpack me dad2y plz3", expected: "Unpack me daddy plzzz"},
	}

	for _, data := range dataValid {
		t.Run(data.name, func(t *testing.T) {
			actual, err := unpack(data.input)

			if err != nil {
				t.Fatalf("Unexpected error: '%s'\n", err.Error())
			}

			if actual != data.expected {
				t.Fatalf("Actual result '%#v' is not equal to expected '%#v'\n", actual, data.expected)
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

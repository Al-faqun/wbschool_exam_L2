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
		{
			name: "Default alpha sort",
			options: SortOptions{},
			input: "MX Linux\nManjaro\nMint\nelementary\nUbuntu",
			expected: "elementary\nManjaro\nMint\nMX Linux\nUbuntu"
		},
		{
			name: "Numeric sort",
			options: SortOptions{},
			input: "1\n10\n2\n21\n23\n3\n432\n5\n5\n60",
			expected: "1\n2\n3\n5\n5\n10\n21\n23\n60\n432" // todo: check how linux sort -n reacts to text 
		},
		{
			name: "Reverse sort",
			options: SortOptions{},
			input: "1.MX Linux\n4.elementary\n2.Manjaro\n5.Ubuntu\n3.Mint",
			expected: "4.elementary\n1.MX Linux\n2.Manjaro\n5.Ubuntu\n3.Mint"
		},
		{
			name: "Remove duplicates & sort",
			options: SortOptions{},
			input: "1.MX Linux\n2.Manjaro\n3.Mint\n4.elementary\n5.Ubuntu\n1.MX Linux\n2.Manjaro\n3.Mint\n4.elementary\n5.Ubuntu",
			expected: "1.MX Linux\n2.Manjaro\n3.Mint\n4.elementary\n5.Ubuntu"
		},
		// {
		// 	name: "Month sort",
		// 	options: SortOptions{}
		// 	input: "March\nFeb\nFebruary\nApril\nAugust\nJuly\nJune\nNovember\nOctober\nDecember\nMay\nSeptember",
		// 	expected: "Jan\nFeb\nFebruary\nMarch\nApril\nMay\nJune\nJuly\nAugust\nSeptember\nOctober\nNovember\nDecember"
		// },
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
package main

import (
	"testing"
)

func TestUnpackValid(t *testing.T) {
	dataValid := []struct {
		name     string
		options  SortOptions
		input    string
		expected string
	}{
		{
			name:     "Default alpha sort",
			options:  SortOptions{},
			input:    "MX Linux\nManjaro\nMint\nelementary\nUbuntu",
			expected: "elementary\nManjaro\nMint\nMX Linux\nUbuntu",
		},
		{
			name:     "Column sort",
			options:  SortOptions{col: 2},
			input:    "Manjaro 300\nArch 100\nMint 400\nNoSecondColumn\nelementary 200\nUbuntu -100",
			expected: "NoSecondColumn\nUbuntu -100\nArch 100\nelementary 200\nManjaro 300\nMint 400",
		},
		{
			name:     "Numeric sort",
			options:  SortOptions{isNum: true},
			input:    "1\n0\n10\n 10\n-10\n10ab\n10ac\n10 0\n100\n21.6789\n elementary\n2\n21.5\n23\n3\n432\nUbuntu\nelementary",
			expected: "-10\n0\n elementary\nelementary\nUbuntu\n1\n2\n3\n 10\n10\n10 0\n10ab\n10ac\n21.5\n21.6789\n23\n100\n432",
		},
		{
			name:     "Reverse sort",
			options:  SortOptions{isRev: true},
			input:    "1.MX Linux\n4.elementary\n2.Manjaro\n5.Ubuntu\n3.Mint",
			expected: "5.Ubuntu\n4.elementary\n3.Mint\n2.Manjaro\n1.MX Linux",
		},
		{
			name:     "Reverse column sort",
			options:  SortOptions{isRev: true, col: 2},
			input:    "Manjaro 300\nMint 400\nNoSecondColumn\nelementary 200\nUbuntu 100",
			expected: "Mint 400\nManjaro 300\nelementary 200\nUbuntu 100\nNoSecondColumn",
		},
		{
			name:     "Remove duplicates & sort",
			options:  SortOptions{isRemDub: true},
			input:    "1.MX Linux\n2.Manjaro\n3.Mint\n4.elementary\n5.Ubuntu\n1.MX Linux\n2.Manjaro\n3.Mint\n4.elementary\n5.Ubuntu",
			expected: "1.MX Linux\n2.Manjaro\n3.Mint\n4.elementary\n5.Ubuntu",
		},
		{
			name:     "Remove duplicates & sort by column",
			options:  SortOptions{isRemDub: true, col: 1},
			input:    "1.MX Linux\n2.Manjaro\n3.Mint\n4.elementary\n5.Ubuntu\n1.MX Linux\n2.Manjaro\n3.Mint\n4.elementary\n5.Ubuntu",
			expected: "1.MX Linux\n2.Manjaro\n3.Mint\n4.elementary\n5.Ubuntu",
		},
	}

	for _, data := range dataValid {
		t.Run(data.name, func(t *testing.T) {
			actual, err := sort(data.input, data.options)

			if err != nil {
				t.Fatalf("Unexpected error: '%q'\n", err.Error())
			}

			if actual != data.expected {
				t.Fatalf("Actual result '%+q' is not equal to expected '%+q'\n", actual, data.expected)
			}
		})
	}
}

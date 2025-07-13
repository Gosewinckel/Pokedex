package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: "  Hello World  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "     ",
			expected: []string{},
		},
		{
			input: "whatyouWouldEXPECT",
			expected: []string{"whatyouwouldexpect"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if  len(c.expected) != len(actual) {
			t.Errorf("actual different length that expected\n actual: %d\n expected: %d\n", len(actual), len(c.expected))
			t.Fail()
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("expected does not match output\n expected: %s\n actual: %s\n", expectedWord, actual)
				t.Fail()
			}
		}
	}
}

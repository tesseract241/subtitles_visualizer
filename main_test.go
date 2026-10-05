package main

import (
	"bufio"
	"strings"
	"testing"
	"unicode/utf8"
)

func Fuzz_scanner_to_channel(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, input string) {
		input = strings.ReplaceAll(strings.ReplaceAll(input, "\n", ""), "\r", "")
		scanner := bufio.NewScanner(strings.NewReader(input))
		channel := _scanner_to_channel(scanner)
		output := ""
		for s := range channel {
			output += s
		}
		if output != input {
			t.Errorf("expected %q, got %q", input, output)
		}
		if utf8.ValidString(input) && !utf8.ValidString(output) {
			t.Errorf("scanner_to_channel produced invalid UTF-8 string %q", output)
		}
	})
}

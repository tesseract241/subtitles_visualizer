package main

import (
	"bufio"
	"strings"
	"testing"
	"unicode/utf8"
)

func Test_parse(t *testing.T) {
	textList := new(TextList)
	_parse("hello", textList)
	if "hello" != textList.ToText() {
		t.Errorf("Expected \"hello\", got %q", textList.ToText())
	}
	_parse(" world", textList)
	if "hello world" != textList.ToText() {
		t.Errorf("Expected \"hello world\", got %q", textList.ToText())
	}
	_parse(" And goodbye", textList)
	if "hello world \nAnd goodbye" != textList.ToText() {
		t.Errorf("Expected \"hello world\"\n And goodbye\", got %q", textList.ToText())
	}
	_parse("Said the old geese", textList)
	if "And goodbye\nSaid the old geese" != textList.ToText() {
		t.Errorf("Expected \"And goodbye\\nSaid the old geese\", got %q", textList.ToText()) 
	}
}

func Fuzz_scanner_to_channel(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, input string) {
		//_scanner_to_channel is meant to strip newlines, so we do it manually on fuzzed inputs
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

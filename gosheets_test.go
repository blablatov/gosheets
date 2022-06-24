package main

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	var tests = []struct {
		spreadsheetId string
		sep           string
	}{
		{"", "\n"},
		{"", ""},
		{"\t", "one\ttwo\tthree\n"},
		{"Data for test", "Number 99999 to data test"},
		{"Yes, no", "No, or, yes"},
	}

	var prevspreadsheetId string
	for _, test := range tests {
		if test.spreadsheetId != prevspreadsheetId {
			fmt.Printf("\n%s\n", test.spreadsheetId)
			prevspreadsheetId = test.spreadsheetId
		}
	}

	var prevsep string
	for _, test := range tests {
		if test.sep != prevsep {
			fmt.Printf("\n%s\n", test.sep)
			prevsep = test.sep
		}
	}
}

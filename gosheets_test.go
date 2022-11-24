package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	var tests = []struct {
		spreadsheetId           string
		sep                     string
		indicator_to_mo_fact_id string
		indicator_to_mo_id      string
		weight                  string
		plan_default            string
		value                   string
	}{
		{"root", "toor\t", "1A663PCe8LUilZ-tWbImbj4vlSikymqRBPA62gDVVddw", "63PCe8LUilZ", "12345", "*)([]{}", "_="},
		{"root", "admin\t", "-tWbImbj4vlSikymqRBPA62gDVVddw=", "--63PCe8LUilZ-", "12345", "*)([]{}", "_="},
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

	var previndicator_to_mo_fact_id string
	for _, test := range tests {
		if test.indicator_to_mo_fact_id != previndicator_to_mo_fact_id {
			fmt.Printf("\n%s\n", test.indicator_to_mo_fact_id)
			previndicator_to_mo_fact_id = test.indicator_to_mo_fact_id
		}
	}

	var previndicator_to_mo_id string
	for _, test := range tests {
		if test.indicator_to_mo_id != previndicator_to_mo_id {
			fmt.Printf("\n%s\n", test.indicator_to_mo_id)
			previndicator_to_mo_id = test.indicator_to_mo_id
		}
	}

	var prevweight string
	for _, test := range tests {
		if test.weight != prevweight {
			fmt.Printf("\n%s\n", test.weight)
			prevweight = test.weight
		}
	}

	var prevplan_default string
	for _, test := range tests {
		if test.plan_default != prevplan_default {
			fmt.Printf("\n%s\n", test.plan_default)
			prevplan_default = test.plan_default
		}
	}

	var prevvalue string
	for _, test := range tests {
		if test.value != prevvalue {
			fmt.Printf("\n%s\n", test.weight)
			prevvalue = test.value
		}
	}
}

func TestSplit(t *testing.T) {
	s, sep := "modbus_tcp:ReadCoils:16", ":"
	durl := strings.Split(s, sep)
	if got, want := len(durl), 3; got != want {
		t.Errorf("Split (%q%q) return %d param, but should: %d", s, sep, got, want)
	}
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	if ctx == nil {
		t.Fatalf("Background returned nil")
	}
	select {
	case x := <-ctx.Done():
		t.Errorf("ctx.Done == %v hould block", x)
	default:
	}
	if gt, wnt := fmt.Sprint(ctx), "context.Background"; gt != wnt {
		t.Errorf("Background().String() = %q want %q", gt, wnt)
	}
}

func checkSize(t *testing.T, path string, size int64) {
	dir, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat %q (looking for size %d): %s", path, size, err)
	}
	if dir.Size() != size {
		t.Errorf("Stat %q: size %d want %d", path, dir.Size(), size)
	}
}

func TestReadFile(t *testing.T) {
	fname := "/.credentials.json"
	fn, err := ioutil.ReadFile(fname)
	if err == nil {
		t.Fatalf("ReadFile %s: error expected, none found", fname)
	}

	fname = "./gosheets_test.go"
	fn, err = ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("ReadFile %s: %v", fname, err)
	}
	checkSize(t, fname, int64(len(fn)))
}

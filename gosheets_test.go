package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
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

func TestReadFile(t *testing.T) {
	fname := "credentials.json"
	fn, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("ReadFile %s: error expected, none found", fname)
	}

	config, err := google.ConfigFromJSON(fn, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		t.Errorf("Не удалось разобрать секретный файл клиента для конфигурации: %v", err)
		client := getClient(config)
		ctx := context.Background()
		_, err := sheets.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			t.Errorf("Не удалось получить Таблицы клиента: %v", err)
		}
	}
}

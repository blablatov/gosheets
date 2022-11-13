package main

import (
	"fmt"
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

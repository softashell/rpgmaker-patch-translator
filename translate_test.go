package main

import "testing"

func TestShouldTranslateText(t *testing.T) {
	type testpair struct {
		text   string
		result bool
	}

	var tests = []testpair{
		{
			`1`,
			false,
		},
		{
			`test`,
			false,
		},
		{
			" ",
			false,
		},
		{
			"あの――",
			true,
		},
	}

	for _, pair := range tests {
		r := shouldTranslateText(pair.text)
		if r != pair.result {
			t.Error(
				"For", pair.text,
				"expected", pair.result,
				"got", r,
			)
		}
	}

}

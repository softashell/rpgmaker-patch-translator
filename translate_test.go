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

func TestCleanTranslatedText(t *testing.T) {
	type testpair struct {
		text   string
		result string
	}

	var tests = []testpair{
		{
			`test`,
			`test`,
		},
		{
			" ",
			" ",
		},
		{
			"あの――",
			"あの――",
		},
		{
			`a good idea of ​​a magician`,
			`a good idea of a magician`,
		},
	}

	for _, pair := range tests {
		r := cleanTranslation(pair.text)
		if r != pair.result {
			t.Errorf("For %q expected %q got %q", pair.text, pair.result, r)
		}
	}

}

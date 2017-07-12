package main

import "testing"

func TestCleanTranslatedText(t *testing.T) {
	type testpair struct {
		input  string
		output string
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
		{
			" ― ―",
			" ー",
		},
		{
			` ー ー ー ー`,
			` ー`,
		},
		{
			`Wow ー っ っ！`,
			`Wow ー！`,
		},
	}

	for _, pair := range tests {
		r := cleanTranslation(pair.input)
		if r != pair.output {
			t.Errorf("For input:\n%q\nexpected:\n%q\ngot:\n%q\n", pair.input, pair.output, r)
		}
	}

}

package text

import (
	"testing"
)

func TestShouldTranslateText(t *testing.T) {
	var tests = []struct {
		input  string
		output bool
	}{
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
		{
			`#####素材アイテム####`,
			true,
		},
		{
			`/\\eS\\[(\\d+),(.*?),(.*?)\\]/`,
			false,
		},
		{
			`/\\<\\s*接触範囲\\s*\\:\\s*(.+?)\\s*\\\>/`,
			false,
		},
	}

	for _, pair := range tests {
		r := ShouldTranslate(pair.input)
		if r != pair.output {
			t.Errorf("For input:\n%q\nexpected: %v got: %v\n", pair.input, pair.output, r)
		}
	}
}

func TestPatchUnescape(t *testing.T) {
	var tests = []struct {
		input  string
		output string
	}{
		{
			`\#\#\#\#\#素材アイテム\#\#\#\#`,
			`#####素材アイテム####`,
		},
		{
			`[守護]水属性ダメージを\\V[20]%軽減`,
			`[守護]水属性ダメージを\V[20]%軽減`,
		},
	}

	for _, pair := range tests {
		r := Unescape(pair.input)
		if r != pair.output {
			t.Errorf("For input:\n%q\nexpected:\n%q\ngot:\n%q\n", pair.input, pair.output, r)
		}
	}
}

func TestPatchEscape(t *testing.T) {
	var tests = []struct {
		input  string
		output string
	}{
		{
			`#####素材アイテム####`,
			`\#\#\#\#\#素材アイテム\#\#\#\#`,
		},
	}

	for _, pair := range tests {
		r := Escape(pair.input)
		if r != pair.output {
			t.Errorf("For input:\n%q\nexpected:\n%q\ngot:\n%q\n", pair.input, pair.output, r)
		}
	}
}

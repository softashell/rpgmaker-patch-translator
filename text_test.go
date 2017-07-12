package main

import (
	"testing"
)

func TestLineBreaking(t *testing.T) {
	var tests = []struct {
		input  string
		output string
	}{
		{
			`\\name[Domestic scent]（…Properly……I want to take a bath……
today…Let's go home.……）`,
			`\\name[Domestic scent]（…Properly……I want to take a bath……
today…Let's go home.……）`,
		},
		{
			`"Money %s\\\\G I got！"`,
			`"Money %s\\\\G I got！"`,
		},
		{
			`☆【A whip】Bamboo that manipulates
thunder。 Get Lightning Lv 20。
attack:＋80 Mausoleum:＋８０ （Blow）（Overall）（Thunder）（Stan）`,
			`☆【A whip】Bamboo that manipulates
thunder。 Get Lightning Lv 20。
attack:＋80 Mausoleum:＋８０ （Blow）（Overall）（
Thunder）（Stan）`,
		},
		{
			`【dagger】It is rusty and its sharpness
【dagger】It is rusty and its sharpness is dull but it will not get stuck。
attack:＋４ （Slashing）（Speed↑）`,
			`【dagger】It is rusty and its sharpness
【dagger】It is rusty and its sharpness is
dull but it will not get stuck。
attack:＋４ （Slashing）（Speed↑）`,
		},
		{
			`【sword】If you equip it, you will get an assault lance song Lv 7。
attack:＋５６ （Slashing）（hit↑）`,
			`【sword】If you equip it, you will get an
assault lance song Lv 7。
attack:＋５６ （Slashing）（hit↑）`,
		},
		{
			`<The contents will be updated even if you see the event by recollection>`,
			`<The contents will be updated even if you
see the event by recollection>`,
		},
		{
			`After viewing the event the content is updated,
Reset by sleeping in bed。
<The contents will be updated even if you see the event by recollection\>`,
			`After viewing the event the content is updated,
Reset by sleeping in bed。
<The contents will be updated even if you
see the event by recollection\>`,
		},
		{
			`During 5 turns, the user is given flames・ice・Give the attribute of lightning`,
			`During 5 turns, the user is given flames・ice・
Give the attribute of lightning`,
		},
	}

	for _, pair := range tests {
		r := breakLines(pair.input)
		if r != pair.output {
			t.Errorf("For input:\n%q\nexpected:\n%q\ngot:\n%q\n", pair.input, pair.output, r)
		}
	}

}

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
		r := shouldTranslateText(pair.input)
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
		r := unescapeText(pair.input)
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
		r := escapeText(pair.input)
		if r != pair.output {
			t.Errorf("For input:\n%q\nexpected:\n%q\ngot:\n%q\n", pair.input, pair.output, r)
		}
	}
}

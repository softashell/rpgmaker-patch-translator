package main

import "testing"

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
		{
			`「For example『Element of fire』When receiving protection of,
　Attack of fire attribute will cause additional attribute attack。
　and,『Element of regeneration』If you receive protection from,
　It will be natural and physical strength will recover during battle」
`,
			`「For example『Element of fire』When
receiving protection of,
　Attack of fire attribute will cause
　additional attribute attack。
　and,『Element of regeneration』If you
　receive protection from,
　It will be natural and physical strength
　will recover during battle」
`,
		},
	}

	lineLength = 42
	lineTolerance = 5

	for _, pair := range tests {
		r := breakLines(pair.input)
		if r != pair.output {
			t.Errorf("For input:\n%q\nexpected:\n%q\ngot:\n%q\n", pair.input, pair.output, r)
		}
	}
}

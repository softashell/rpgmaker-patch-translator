package main

import "testing"

func TestLineBreaking(t *testing.T) {
	var tests = []struct {
		text   string
		result string
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
	}

	for _, pair := range tests {
		r := breakLines(pair.text)
		if r != pair.result {
			t.Errorf("For\n%q\nexpected\n%q\ngot\n%q\n", pair.text, pair.result, r)
		}
	}

}

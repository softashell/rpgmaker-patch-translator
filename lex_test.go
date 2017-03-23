package main

import "testing"

type testpair struct {
	value  string
	result string
}

func TestRawTextExtraction(t *testing.T) {
	var tests = []testpair{
		{
			`An undocumented(test)`,
			`An undocumented(test)`,
		},
		{
			`An undocumented `,
			`An undocumented `,
		},
		{
			`Adventurer's clothes`,
			`Adventurer's clothes`,
		},
		{
			`An undocumented if`,
			`An undocumented if`,
		},
		{
			`An undocumented if()`,
			`An undocumented `,
		},
		{
			`if(v[178] \\>= 40)`,
			``,
		},
		{
			`踊れ if(v[178] \\>= 40)`,
			`踊れ `,
		},
		{
			`\\\>\\C[14]…今は使用できません。`,
			`…今は使用できません。`,
		},
		{
			`Basic Switch \u0026 Variable`,
			`Basic Switch  Variable`,
		},
	}

	for _, pair := range tests {
		v := getOnlyText(pair.value)
		if v != pair.result {
			t.Error(
				"For", pair.value,
				"expected", pair.result,
				"got", v,
			)
		}
	}

}

func TestTranslatableTextExtraction(t *testing.T) {
	var tests = []testpair{
		{
			`if(v[178] \\>= 40)`,
			``,
		},
		{
			`踊れ if(v[178] \\>= 40)`,
			`踊れ `,
		},
		{
			`\\\>\\C[14]…今は使用できません。`,
			`今は使用できません`,
		},
		{
			`\\\>…「\\C[14]\\N[2]\\C[0]」は、\\\>　カナーン村の酒場へ帰っていった。`,
			`は、カナーン村の酒場へ帰っていった`,
		},
		{
			`\\\>…「\\C[14]\\N[2]\\C[0]」の同伴時間が\\C[3] 10 \\C[0]を超えていますので、\\\>　「\\C[14]\\N[2]\\C[0]」の欲情度が\\C[3] 100 \\C[0]未満の場合、\\\>　時間経過で上昇していきます。`,
			`の同伴時間が 10 を超えていますので、の欲情度が 100 未満の場合、時間経過で上昇していきます`,
		},
		{
			`%sを %s 回復した！`,
			`を  回復した！`,
		},
		{
			`"お金を %s\\\\G 手に入れた！"`,
			``,
		},
		{
			`'\\.'`,
			`''`,
		},
		{
			`【\\C[14]\\N[2]\\C[0]】　\\{アハァァーーーンッ！！`,
			`アハァァーーーンッ！！`,
		},
		{
			`\\\>…「\\C[10]バッポウ\\C[0]」再出現カウント： \\C[3] \\V[491] \\C[0]\\\>\\\>　\\C[14]※出現場所：キータニ平原\\\>　\\C[14]※カウントが 0 になると再挑戦可能になります。`,
			`バッポウ再出現カウント：   ※出現場所：キータニ平原※カウントが 0 になると再挑戦可能になります`,
		},
		{
			`【\\C[14]骨董屋\\C[0]】　「謎の防具」か。　ほぅ、複数持っているようだな。　一気に鑑定するかい？\\C[3] \\V[982] G\\C[0] 頂くけどな。\\$`,
			`骨董屋謎の防具かほぅ、複数持っているようだな一気に鑑定するかい？  G 頂くけどな`,
		},
		{
			`【\\C[14]ザウナー\\C[0]】　\\}…待て待て、この場を乗り切るための詭弁さ。\\{　\\}すまないが、我慢して様子を見ててくれ。\\{　\\}必ず上手くいくさ。`,
			`ザウナー待て待て、この場を乗り切るための詭弁さすまないが、我慢して様子を見ててくれ必ず上手くいくさ`,
		},
		{
			`\\c[3]宿屋の主人\\c[0]ほほ～、シスターとは珍しい！\\lこんな辺鄙な島で布教活動かね？言っとくが、ワシは神など信じないぞ。`,
			`宿屋の主人ほほ、シスターとは珍しい！こんな辺鄙な島で布教活動かね？言っとくが、ワシは神など信じないぞ`,
		},
		{
			`0\\G 手に入れた！`,
			`0 手に入れた！`,
		},
		{
			`\\\>牡丹の命が５回復した。\\|\\.\\^`,
			`牡丹の命が５回復した`,
		},
		{
			`氷結水　(残\\V[29]) en(s[28])`,
			`氷結水残 `,
		},
		{
			`疾風苦無(消費1)　en(v[25] \>= 1)`,
			`疾風苦無消費1`,
		},
	}

	for _, pair := range tests {
		items := parseText(pair.value)

		var v string

		for _, item := range items {
			if item.typ == itemText {
				v += item.val
			}
		}

		if v != pair.result {
			t.Error(
				"For", pair.value,
				"expected", pair.result,
				"got", v,
			)
		}
	}
}

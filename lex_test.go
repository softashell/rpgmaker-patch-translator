package main

import (
	"testing"

	log "github.com/Sirupsen/logrus"
)

type testpair struct {
	value  string
	result string
}

func TestTextExtraction(t *testing.T) {
	log.SetLevel(log.DebugLevel)

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

func TestTextTranslation(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	var tests = []testpair{
		{
			`if(v[178] \\>= 40)`,
			`if(v[178] \\>= 40)`,
		},
		{
			`踊れ if(v[178] \\>= 40)`,
			`Dance if(v[178] \\>= 40)`,
		},
		{
			`\\\>\\C[14]…今は使用できません。`,
			`\\\>\\C[14]…It can not be used now。`,
		},
		{
			`\\\>…「\\C[14]\\N[2]\\C[0]」は、\\\>　カナーン村の酒場へ帰っていった。`,
			`\\\>…「\\C[14]\\N[2]\\C[0]」,\\\>　I came back to the bar in Kanan village。`,
		},
		{
			`\\\>…「\\C[14]\\N[2]\\C[0]」の同伴時間が\\C[3] 10 \\C[0]を超えていますので、\\\>　「\\C[14]\\N[2]\\C[0]」の欲情度が\\C[3] 100 \\C[0]未満の場合、\\\>　時間経過で上昇していきます。`,
			`\\\>…「\\C[14]\\N[2]\\C[0]」The accompanying time of\\C[3] Ten \\C[0]Because it exceeds,\\\>　「\\C[14]\\N[2]\\C[0]」The degree of desire is\\C[3] 100 \\C[0],\\\>　It will rise over time。`,
		},
		{
			`%sを %s 回復した！`,
			`%sTo %s He recovered!`,
		},
		{
			`"お金を %s\\\\G 手に入れた！"`,
			`"お金を %s\\\\G 手に入れた！"`,
		},
		{
			`'\\.'`,
			`'\\.'`,
		},
		{
			`【\\C[14]\\N[2]\\C[0]】　\\{アハァァーーーンッ！！`,
			`【\\C[14]\\N[2]\\C[0]】　\\{Ahhhhhhhhhh! It is!`,
		},
		{
			`\\\>…「\\C[10]バッポウ\\C[0]」再出現カウント： \\C[3] \\V[491] \\C[0]\\\>\\\>　\\C[14]※出現場所：キータニ平原\\\>　\\C[14]※カウントが 0 になると再挑戦可能になります。`,
			`\\\>…「\\C[10]Bapau\\C[0]」Reappearance count: \\C[3] \\V[491] \\C[0]\\\>\\\>　\\C[14]※ Appearance place: Ketani Plain\\\>　\\C[14]※ It will be possible to challenge again when the count reaches 0。`,
		},
		{
			`【\\C[14]骨董屋\\C[0]】　「謎の防具」か。　ほぅ、複数持っているようだな。　一気に鑑定するかい？\\C[3] \\V[982] G\\C[0] 頂くけどな。\\$`,
			`【\\C[14]Antique shop\\C[0]】　「Mystery armor」Or。　Huh, it seems to have more than one.。　Do you want to appreciate it at once?\\C[3] \\V[982] G\\C[0] I will get it。\\$`,
		},
		{
			`【\\C[14]ザウナー\\C[0]】　\\}…待て待て、この場を乗り切るための詭弁さ。\\{　\\}すまないが、我慢して様子を見ててくれ。\\{　\\}必ず上手くいくさ。`,
			`【\\C[14]Sauna\\C[0]】　\\}…Wait, wait a moment to survive this place。\\{　\\}Sorry, please be patient and look at the situation。\\{　\\}It certainly will be fine.。`,
		},
		{
			`\\c[3]宿屋の主人\\c[0]ほほ～、シスターとは珍しい！\\lこんな辺鄙な島で布教活動かね？言っとくが、ワシは神など信じないぞ。`,
			`\\c[3]Innkeeper\\c[0]Cheeks～, Sister and unusual!\\lIs this missionary activity on such a remote island? I will tell you, Eagle does not believe in God。`,
		},
		{
			`0\\G 手に入れた！`,
			`0\\G I got!`,
		},
		{
			`\\\>牡丹の命が５回復した。\\|\\.\\^`,
			`\\\>Peony's life recovered 5。\\|\\.\\^`,
		},
		{
			`氷結水　(残\\V[29]) en(s[28])`,
			`Iced water　(Remaining\\V[29]) en(s[28])`,
		},
		{
			`疾風苦無(消費1)　en(v[25] \>= 1)`,
			`Shippuuden (Consumption 1)　en(v[25] \>= 1)`,
		},
	}

	for _, pair := range tests {
		items := parseText(pair.value)
		v := translateItems(items)
		if v != pair.result {
			t.Error(
				"For", pair.value,
				"expected", pair.result,
				"got", v,
			)
		}
	}
}

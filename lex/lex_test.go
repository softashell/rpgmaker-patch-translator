package lex

import (
	"testing"

	"gitgud.io/softashell/rpgmaker-patch-translator/text"
	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
)

type testpair struct {
	input  string
	output string
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
		{
			`Mun「Ha～～～\\\\!`,
			`Mun「Ha～～～`,
		},
		{
			`Marcus「Hey there.……！\\\\!`,
			`Marcus「Hey there.……！`,
		},
		{
			`Yorkie「"Suddenly it is a question, I have been worried since long ago`,
			`Yorkie「"Suddenly it is a question, I have been worried since long ago`,
		},
		{
			`"0x#{text}"`,
			`"0x"`,
		},
		{
			`/<\#{GRPLUS::M_WORD}[：:](\\S+)\>/`,
			`/<[：:](+)>/`,
		},
	}

	for _, pair := range tests {
		r := getOnlyText(text.Unescape(pair.input))
		if r != pair.output {
			t.Errorf("For input:\n%q\nexpected:\n%q\ngot:\n%q\n", pair.input, pair.output, r)
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
			`の同伴時間がを超えていますので、の欲情度が未満の場合、時間経過で上昇していきます`,
		},
		{
			`%sを %s 回復した！`,
			`を  回復した`,
		},
		{
			`"お金を %s\\\\G 手に入れた！"`,
			`お金を  手に入れた`,
		},
		{
			`'\\.'`,
			``,
		},
		{
			`【\\C[14]\\N[2]\\C[0]】　\\{アハァァーーーンッ！！`,
			`アハァァーーーンッ`,
		},
		{
			`\\\>…「\\C[10]バッポウ\\C[0]」再出現カウント： \\C[3] \\V[491] \\C[0]\\\>\\\>　\\C[14]※出現場所：キータニ平原\\\>　\\C[14]※カウントが 0 になると再挑戦可能になります。`,
			`バッポウ再出現カウント※出現場所キータニ平原※カウントが 0 になると再挑戦可能になります`,
		},
		{
			`【\\C[14]骨董屋\\C[0]】　「謎の防具」か。　ほぅ、複数持っているようだな。　一気に鑑定するかい？\\C[3] \\V[982] G\\C[0] 頂くけどな。\\$`,
			`骨董屋謎の防具かほぅ、複数持っているようだな一気に鑑定するかい 頂くけどな`,
		},
		{
			`【\\C[14]ザウナー\\C[0]】　\\}…待て待て、この場を乗り切るための詭弁さ。\\{　\\}すまないが、我慢して様子を見ててくれ。\\{　\\}必ず上手くいくさ。`,
			`ザウナー待て待て、この場を乗り切るための詭弁さすまないが、我慢して様子を見ててくれ必ず上手くいくさ`,
		},
		{
			`\\c[3]宿屋の主人\\c[0]ほほ～、シスターとは珍しい！\\lこんな辺鄙な島で布教活動かね？言っとくが、ワシは神など信じないぞ。`,
			`宿屋の主人ほほ、シスターとは珍しいこんな辺鄙な島で布教活動かね言っとくが、ワシは神など信じないぞ`,
		},
		{
			`0\\G 手に入れた！`,
			` 手に入れた`,
		},
		{
			`\\\>牡丹の命が５回復した。\\|\\.\\^`,
			`牡丹の命が５回復した`,
		},
		{
			`氷結水　(残\\V[29]) en(s[28])`,
			`氷結水残`,
		},
		{
			`疾風苦無(消費1)　en(v[25] \>= 1)`,
			`疾風苦無消費1`,
		},
		{
			`\\\\B\\\\I\\\\C[4]レジネッタ：\\\\C[0]\\\\/I\\\\/B\nあ……ふぁっ……！\n`,
			`レジネッタあふぁっ`,
		},
		{
			`\\\\Bポータルフリントを手に入れた！`,
			`ポータルフリントを手に入れた`,
		},
		{
			`\\name[優理香]懐かしいなぁ。`,
			`優理香懐かしいなぁ`,
		},
		{
			`PT加入en(!s[484] and v[25] <2)`,
			`PT加入`,
		},
		{
			`#####素材アイテム####`,
			`素材アイテム`,
		},
		{
			`\\i[21]メンタルキュア`,
			`メンタルキュア`,
		},
		{
			`<< 迷宮入口へ \>\>`,
			` 迷宮入口へ `,
		},
		{
			`@3「あらら、今は入れないのかぁ。仕方ないな……また改めて来よう」`,
			`あらら、今は入れないのかぁ仕方ないなまた改めて来よう`,
		},
		{
			`@-1「あらら」@15`,
			`あらら`,
		},
		{
			`/(?:付加ポップアップ非表示|add_no_display)/`,
			`付加ポップアップ非表示`,
		},
		{
			`/(?:ポップアップ表示名|display_name)\s*=\s*"([^"]*)"/`,
			`ポップアップ表示名`,
		},
		{
			`/<#{S_B_D::N}[:：](\S+),(\S+)>/`,
			``,
		},
		{
			`"LNX11a:バトラーグラフィック指定の引数が正しくありません。"`,
			`バトラーグラフィック指定の引数が正しくありません`,
		},
		{
			`/(?:解除ポップアップ表示名|remove_display_name)\s*=\s*"([^"]*)"/`,
			`解除ポップアップ表示名`,
		},
		{
			`"OK:LNX11b_リフォーム・バトルステータス"`,
			`リフォームバトルステータス`,
		},
		{
			`\\1	こんにちわ、シスター。`,
			`こんにちわ、シスター`,
		},
		{
			`[レース10]`,
			`レース10`,
		},
		{
			`---以下練成アイテム`,
			`以下練成アイテム`,
		},
		{
			`---------シーフスキルリスト`,
			`シーフスキルリスト`,
		},
		{
			`[武器]自身が使用するデバフの効果量が\\V[19]%UP`,
			`武器自身が使用するデバフの効果量が`,
		},
	}

	for _, pair := range tests {
		items, err := ParseText(text.Unescape(pair.input))
		if err != nil {
			log.Errorf("%s\ntext: %q", err, pair.input)
			log.Error(spew.Sdump(items))
		} else {
			log.Debug(spew.Sdump(items))
		}

		var r string

		for _, item := range items {
			if (item.Typ == ItemText || item.Typ == ItemNumber) && text.ShouldTranslate(item.Val) {
				r += item.Val
			}
		}

		if r != pair.output {
			t.Errorf("For input:\n%q\nexpected:\n%q\ngot:\n%q\n", pair.input, pair.output, r)
		}
	}
}

/*
func TestNumberExtraction(t *testing.T) {
	var tests = []testpair{
		{
			`「そ、そうなのかい……、さ、30万かぁー」`,
			`30万`,
		},
		{
			`0\G 手に入れた！`,
			`0`,
		},
		{
			`牡丹の命が５回復した。\|\.\^`,
			`５`,
		},
		{
			`@-1「あらら」@15`,
			``,
		},
		{
			`"0x#{text}"`,
			`0`,
		},
		{
			`[レース10]`,
			`10`,
		},
		{
			`ＨＰとＭＰを１００％回復する`,
			`１００％`,
		},
		{
			`でも３階層のモンスターはやたら攻撃力が高いからな。`,
			`３階`,
		},
		{
			`２階層の敵は魔法防御が低いものが多い。`,
			`２階`,
		},
	}

	for _, pair := range tests {
		items, err := ParseText(pair.input)

		if err != nil {
			log.Errorf("%s\ntext: %q", err, pair.input)
			log.Error(spew.Sdump(items))
		} else {
			log.Debug(spew.Sdump(items))
		}

		var r string

		for _, item := range items {
			if item.Typ == ItemNumber {
				r += item.Val
			}
		}

		if r != pair.output {
			t.Errorf("For input:\n%q\nexpected:\n%q\ngot:\n%q\nitems:\n%s", pair.input, pair.output, r, spew.Sdump(items))
		}
	}
}
*/

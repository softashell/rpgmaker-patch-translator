package lex

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func TestItemAssembly(t *testing.T) {
	var tests = []testpair{
		{
			`>\C[14]…今は使用できません。`,
			`>\C[14]…今は使用できません。`,
		},
		{
			`【\C[14]\N[2]\C[0]】　\{アハァァーーーンッ！！`,
			`【\C[14]\N[2]\C[0]】　\{ アハァァーーーンッ！！`,
		},
		{
			`疾風苦無(消費1)　en(v[25] >= 1)`,
			`疾風苦無 (消費1)　en(v[25] >= 1)`,
		},
		{
			`PT加入en(!s[484] and v[25] <2)`,
			`PT加入 en(!s[484] and v[25] <2)`,
		},
		{
			`\i[21]メンタルキュア`,
			`\i[21] メンタルキュア`,
		},
		{
			`[レース10]`,
			`[レース10]`,
		},
		{
			`%sの%sを %s 奪った`,
			`%s の %s を %s 奪った`,
		},
		{
			`[武器]攻撃時に\V[19]%の闇属性追加ダメージ`,
			`[武器]攻撃時に \V[19]%の闇属性追加ダメージ`,
		},
	}

	for _, pair := range tests {
		//pair.input = unescapeText(pair.input)
		items, err := ParseText(pair.input)
		if err != nil {
			log.Errorf("%s\ntext: %q", err, pair.input)
			log.Error(spew.Sdump(items))
		} else {
			log.Debug(spew.Sdump(items))
		}

		r := assembleItems(items)

		if r != pair.output {
			t.Errorf("For input:\n%q\nexpected:\n%q\ngot:\n%q\nitems:\n%s", pair.input, pair.output, r, spew.Sdump(items))
		}
	}
}

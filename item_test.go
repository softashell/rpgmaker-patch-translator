package main

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
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
			`疾風苦無 (消費 1)　en(v[25] >= 1)`,
		},
		{
			`PT加入en(!s[484] and v[25] <2)`,
			`PT 加入 en(!s[484] and v[25] <2)`,
		},
		{
			`\i[21]メンタルキュア`,
			`\i[21] メンタルキュア`,
		},
		{
			`"0x#{text}"`,
			`"0x#{text}"`,
		},
		{
			`[レース10]`,
			`[レース 10]`,
		},
	}

	for _, pair := range tests {
		//pair.input = unescapeText(pair.input)
		items, err := parseText(pair.input)
		if err != nil {
			log.Errorf("%s\ntext: %q", err, pair.input)
			log.Error(spew.Sdump(items))
		} else {
			log.Debug(spew.Sdump(items))
		}

		var r string

		// Use orignal text as translation
		for i := range items {
			if items[i].typ == itemText {
				items[i].trans += items[i].val
			}
		}

		r = assembleItems(items)

		if r != pair.output {
			t.Errorf("For input:\n%q\nexpected:\n%q\ngot:\n%q\n", pair.input, pair.output, r)
		}
	}
}

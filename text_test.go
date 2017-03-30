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
	}

	for _, pair := range tests {
		r := breakLines(pair.text)
		if r != pair.result {
			t.Errorf("For\n%q\nexpected\n%q\ngot\n%q\n", pair.text, pair.result, r)
		}
	}

}

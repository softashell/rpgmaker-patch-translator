package main

import (
	"testing"
)

func Test_shouldTranslateContextVX(t *testing.T) {
	type args struct {
		c    string
		text string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"InlineScript",
			args{
				`> CONTEXT: Commonevents/50/4/InlineScript/1:12`,
				`"メニュー背景"`,
			},
			false,
		},
		{
			"Actor name",
			args{
				`: Actors/1/name/`,
				``,
			},
			true,
		},
		{
			"Armor name",
			args{
				`: Armors/1/name/`,
				``,
			},
			true,
		},
		{
			"Item name",
			args{
				`: Items/1/name/`,
				``,
			},
			true,
		},
		{
			"Class name",
			args{
				`: Classes/1/name/`,
				``,
			},
			true,
		},
		{
			"Window_NameInput",
			args{
				`: Scripts/Window_NameInput/46:36`,
				``,
			},
			false,
		},
		{
			"Window_Message",
			args{
				`: Scripts/Window_Message/357:10`,
				``,
			},
			false,
		},
		{
			"Currency",
			args{
				`: System/currency_unit/`,
				`万円`,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldTranslateContextVX(tt.args.c, tt.args.text); got != tt.want {
				t.Errorf("shouldTranslateContextVX() = %v, want %v", got, tt.want)
			}
		})
	}
}

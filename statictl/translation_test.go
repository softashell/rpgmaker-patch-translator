package statictl

import (
	"regexp"
	"testing"
)

func TestDb_GetDynamic(t *testing.T) {
	type fields struct {
		dbRe translationDBRegexMap
	}
	type args struct {
		str string
		typ TranslationType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Benis",
			fields: fields{
				dbRe: translationDBRegexMap{
					TransGeneric: []translationDBRegex{{
						regex:       regexp.MustCompile(`(b(en))is`),
						replacement: "$1$1$2",
					}},
				},
			},
			args: args{
				str: "benis bepis",
				typ: TransGeneric,
			},
			want:    "benbenen bepis",
			wantErr: false,
		},
		{
			name: "Benis no match",
			fields: fields{
				dbRe: translationDBRegexMap{
					TransGeneric: []translationDBRegex{{
						regex:       regexp.MustCompile(`(ben)is`),
						replacement: "$1$1",
					}},
				},
			},
			args: args{
				str: "ben",
				typ: TransGeneric,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl := &Db{
				dbRe: tt.fields.dbRe,
			}
			got, err := tl.getDynamic(tt.args.str, tt.args.typ)
			if (err != nil) != tt.wantErr {
				t.Errorf("Db.GetDynamic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Db.GetDynamic() = %v, want %v", got, tt.want)
			}
		})
	}
}

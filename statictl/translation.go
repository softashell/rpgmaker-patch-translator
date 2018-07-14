package statictl

import (
	"fmt"
	"strings"
)

func (t *Db) GetTranslation(str string, typ TranslationType) (string, error) {
	// TODO: Maybe only trim right side?
	str = strings.TrimSpace(str)

	// Full replacement
	tl, err := t.GetStatic(str, typ)
	if err == nil {
		return tl, err
	}

	// Regex search and replace
	tl, err = t.GetDynamic(str, typ)
	if err == nil {
		return tl, err
	}

	// Fall back to generic
	if typ != TransGeneric {
		tl, err = t.GetStatic(str, TransGeneric)
		if err == nil {
			return tl, err
		}

		tl, err = t.GetDynamic(str, TransGeneric)
		if err == nil {
			return tl, err
		}
	}

	return "", fmt.Errorf("no translation")
}

func (t *Db) GetStatic(str string, typ TranslationType) (string, error) {
	if tl, ok := t.db[typ][str]; ok {
		return tl, nil
	}

	return "", fmt.Errorf("no translation")
}

func (t *Db) GetDynamic(str string, typ TranslationType) (string, error) {
	var match bool

	for _, r := range t.dbRe[typ] {
		if !r.regex.MatchString(str) {
			continue
		}

		match = true

		str = r.regex.ReplaceAllString(str, r.replacement)
	}

	if match {
		return str, nil
	}

	return "", fmt.Errorf("no translation")
}

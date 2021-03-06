package statictl

import (
	"fmt"
	"strings"
)

// GetTranslation Replaces source if match is found and stops any further translation.
// Called after RunPreTranslation
func (t *Db) GetTranslation(str string, typ TranslationType) (string, error) {
	// TODO: Maybe only trim right side?
	str = strings.TrimSpace(str)

	// Full replacement
	tl, err := t.getStatic(str, typ)
	if err == nil {
		return tl, err
	}

	// Regex search and replace
	tl, err = t.getDynamic(str, typ)
	if err == nil {
		return tl, err
	}

	// Fall back to generic
	if typ != TransGeneric {
		tl, err = t.getStatic(str, TransGeneric)
		if err == nil {
			return tl, err
		}

		tl, err = t.getDynamic(str, TransGeneric)
		if err == nil {
			return tl, err
		}
	}

	return "", fmt.Errorf("no translation")
}

func (t *Db) getStatic(str string, typ TranslationType) (string, error) {
	if tl, ok := t.db[typ][str]; ok {
		return tl, nil
	}

	return "", fmt.Errorf("no translation")
}

func (t *Db) getDynamic(str string, typ TranslationType) (string, error) {
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

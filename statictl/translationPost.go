package statictl

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

// RunPostTranslation Edits text returned from translation service, ignores static tl
func (t *Db) RunPostTranslation(str string) (string, error) {
	str = strings.TrimSpace(str) // Might not be a good idea

	var err error

	str, err = t.applyPostStatic(str, TransGeneric)
	if err != nil {
		log.Error(err)
	}

	str, err = t.applyPostDynamic(str, TransGeneric)
	if err != nil {
		log.Error(err)
	}

	return str, nil
}

func (t *Db) applyPostStatic(str string, typ TranslationType) (string, error) {
	for strFind, strSub := range t.dbPost[typ] {
		str = strings.Replace(str, strFind, strSub, -1)
	}

	return str, nil
}

func (t *Db) applyPostDynamic(str string, typ TranslationType) (string, error) {
	for _, r := range t.dbRePost[typ] {
		if !r.regex.MatchString(str) {
			continue
		}

		str = r.regex.ReplaceAllString(str, r.replacement)
	}

	return str, nil
}

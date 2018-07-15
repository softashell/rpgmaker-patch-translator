package statictl

import (
	"strings"

	log "github.com/Sirupsen/logrus"
)

// RunPostTranslation Edits text returned from translation service, ignores static tl
func (t *Db) RunPostTranslation(str string) (string, error) {
	str = strings.TrimSpace(str) // Might not be a good idea

	var err error

	str, err = t.ApplyPreStatic(str, TransGeneric)
	if err != nil {
		log.Error(err)
	}

	str, err = t.ApplyPreDynamic(str, TransGeneric)
	if err != nil {
		log.Error(err)
	}

	return str, nil
}

func (t *Db) ApplyPostStatic(str string, typ TranslationType) (string, error) {
	for strFind, strSub := range t.dbPost[typ] {
		str = strings.Replace(str, strFind, strSub, -1)
	}

	return str, nil
}

func (t *Db) ApplyPostDynamic(str string, typ TranslationType) (string, error) {
	for _, r := range t.dbRePost[typ] {
		if !r.regex.MatchString(str) {
			continue
		}

		str = r.regex.ReplaceAllString(str, r.replacement)
	}

	return str, nil
}

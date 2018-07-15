package statictl

import (
	"strings"

	log "github.com/Sirupsen/logrus"
)

// RunPreTranslation Edits original text, called before anything else
func (t *Db) RunPreTranslation(str string, typ TranslationType) (string, error) {
	str = strings.TrimSpace(str)

	var err error

	str, err = t.ApplyPreStatic(str, typ)
	if err != nil {
		log.Error(err)
	}

	str, err = t.ApplyPreDynamic(str, typ)
	if err != nil {
		log.Error(err)
	}

	return str, nil
}

func (t *Db) ApplyPreStatic(str string, typ TranslationType) (string, error) {
	for strFind, strSub := range t.dbPre[typ] {
		str = strings.Replace(str, strFind, strSub, -1)
	}

	return str, nil
}

func (t *Db) ApplyPreDynamic(str string, typ TranslationType) (string, error) {
	for _, r := range t.dbRePre[typ] {
		if !r.regex.MatchString(str) {
			continue
		}

		str = r.regex.ReplaceAllString(str, r.replacement)
	}

	return str, nil
}

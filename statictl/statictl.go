package statictl

import (
	"log"
	"regexp"
)

type translationDBMap map[TranslationType]translationDB
type translationDB map[string]string

type translationDBRegexMap map[TranslationType][]translationDBRegex
type translationDBRegex struct {
	regex       *regexp.Regexp
	replacement string
}

type Db struct {
	db   translationDBMap
	dbRe translationDBRegexMap
}

func New() *Db {
	t := &Db{}

	t.db = make(translationDBMap)
	t.dbRe = make(translationDBRegexMap)

	err := t.loadDatabases()
	if err != nil {
		log.Fatal("Failed to load static translations")
	}

	return t
}

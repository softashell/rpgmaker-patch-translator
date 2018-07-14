package statictl

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/hjson/hjson-go"
)

// Normal replace
type DatabaseEntryStatic struct {
	Original   string
	Translated string
}

// Regex
type DatabaseEntryDynamic struct {
	RegexMatch   string
	RegexReplace string

	/*
		re := regexp.MustCompile(`((fo)(o))`)
		s := re.ReplaceAllString("foo", "$1-$1-$2-$3") // foo => foo-foo-fo-o
	*/
}

var staticDBPath = filepath.Join("database", "static")
var dynamicDBPath = filepath.Join("database", "dynamic")

func (t *Db) loadDatabases() error {
	for k, v := range databaseFiles {
		err := t.loadDatabase(k, v)
		if err != nil {
			continue
		}
	}

	return nil
}

func (t *Db) loadDatabase(tlType TranslationType, fileName string) error {
	filePathStatic := filepath.Join(staticDBPath, fileName)
	filePathDynamic := filepath.Join(dynamicDBPath, fileName)

	var db translationDB
	var dbRe []translationDBRegex
	var err error

	if !fileExists(filePathStatic) {
		log.Infof("Database %s doesn't exist, creating empty file", filePathStatic)
		createEmptyDatabase(filePathStatic, DbStatic)
	} else {
		db, err = t.loadDatabaseStatic(filePathStatic)
		if err != nil {
			log.Error(err)
		}
	}

	if !fileExists(filePathDynamic) {
		log.Infof("Database %s doesn't exist, creating empty file", filePathDynamic)
		createEmptyDatabase(filePathDynamic, DbDynamic)
	} else {
		dbRe, err = t.loadDatabaseDynamic(filePathDynamic)
		if err != nil {
			log.Error(err)
		}
	}

	t.db[tlType] = db
	t.dbRe[tlType] = dbRe

	return nil
}

func (t *Db) loadDatabaseStatic(fileName string) (translationDB, error) {
	db := make(translationDB)

	log.Infof("Parsing database %s", fileName)

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("File error: %v\n", err)
	}

	var dat []interface{}
	if err := hjson.Unmarshal(file, &dat); err != nil {
		log.Error(err)
	} else {
		// convert to JSON
		b, _ := json.Marshal(dat)

		// unmarshal
		var dbEntries []DatabaseEntryStatic
		json.Unmarshal(b, &dbEntries)

		for _, v := range dbEntries {
			db[v.Original] = v.Translated
		}
	}

	return db, nil
}

func (t *Db) loadDatabaseDynamic(fileName string) ([]translationDBRegex, error) {
	var dbRe []translationDBRegex

	log.Debugf("Parsing database %s", fileName)

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("File error: %v\n", err)
	}

	var dat []interface{}
	if err := hjson.Unmarshal(file, &dat); err != nil {
		log.Error(err)
	} else {
		// convert to JSON
		b, _ := json.Marshal(dat)

		// unmarshal
		var dbEntries []DatabaseEntryDynamic
		json.Unmarshal(b, &dbEntries)

		for _, v := range dbEntries {
			if len(v.RegexMatch) < 1 {
				continue
			}

			matchRegex, err := regexp.Compile(v.RegexMatch)
			if err != nil {
				log.Errorf("Failed to compile regex\n%s", v.RegexMatch)
				continue
			}

			dbRe = append(dbRe, translationDBRegex{
				regex:       matchRegex,
				replacement: v.RegexReplace,
			})
		}
	}

	return dbRe, nil
}

func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	// Probably exists but it might have other issues
	return true
}

func createEmptyDatabase(filePath string, typ DatabaseType) error {

	var content interface{}
	var fileDir string

	switch typ {
	case DbStatic:
		fileDir = staticDBPath
		content = []DatabaseEntryStatic{{
			Original:   "",
			Translated: "",
		}}

	case DbDynamic:
		fileDir = dynamicDBPath
		content = []DatabaseEntryDynamic{{
			RegexMatch:   "",
			RegexReplace: "",
		}}
	}

	err := os.MkdirAll(fileDir, 0755)
	if err != nil && !os.IsExist(err) {
		log.Error(err)
		return err
	}

	outJSON, err := hjson.Marshal(content)
	if err != nil {
		log.Error(err)
	}

	err = ioutil.WriteFile(filePath, outJSON, 0644)

	return err
}

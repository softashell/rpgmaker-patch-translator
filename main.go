package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

func main() {
	//	log.SetLevel(log.DebugLevel)

	args := os.Args
	if len(args) < 2 {
		log.Error("Program requires patch directory as argument")
		os.Exit(1)
	}

	dir := os.Args[1]

	err := checkPatchVersion(dir)
	check(err)

	fileList := getDirectoryContents(path.Join(dir, "Patch"))
	if len(fileList) < 1 {
		log.Error("Couldn't find anything to translate")
		os.Exit(1)
	}

	for _, file := range fileList {
		//		go func(file string) {
		patch, err := parsePatchFile(file)
		check(err)

		patch, err = translatePatch(patch)
		check(err)

		err = writePatchFile(patch)
		check(err)

		log.Info("Translated ", file)
		//		}(file)

		//var input string
		//fmt.Scanln(&input)
	}
}

func checkPatchVersion(dir string) error {
	contents, err := ioutil.ReadFile(path.Join(dir, "RPGMKTRANSPATCH"))
	if err != nil {
		return err
	}

	if string(contents) != "> RPGMAKER TRANS PATCH V3" {
		err = fmt.Errorf("Unsupported patch version")
	}

	return err
}

func getDirectoryContents(dir string) []string {
	fileList := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		fileList = append(fileList, path)

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return fileList
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

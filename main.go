package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
)

func main() {
	// TODO: Add flags
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Program requires patch directory as argument")
	}

	dir := os.Args[1]

	err := checkPatchVersion(dir)
	if err != nil {
		log.Fatal(err)
	}

	fileList := getDirectoryContents(filepath.Join(dir, "Patch"))
	if len(fileList) < 1 {
		log.Fatal("Couldn't find anything to translate")
	}

	count := len(fileList)
	start := time.Now()

	for i, file := range fileList {
		fmt.Printf("Processing %q (%d/%d)\n", filepath.Base(file), i+1, count)

		patch, err := parsePatchFile(file)
		if err != nil {
			log.Error(err)
			continue
		}

		patch, err = translatePatch(patch)
		check(err)

		err = writePatchFile(patch)
		check(err)
	}

	fmt.Printf("Finished in %s\n", time.Since(start))
}

func checkPatchVersion(dir string) error {
	contents, err := ioutil.ReadFile(filepath.Join(dir, "RPGMKTRANSPATCH"))
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
		if f.IsDir() || filepath.Ext(path) != ".txt" {
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

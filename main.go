package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

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

	count := len(fileList)
	start := time.Now()

	for i, file := range fileList {
		fmt.Printf("Processing %q (%d/%d)\n", path.Base(file), i+1, count)

		patch, err := parsePatchFile(file)
		check(err)

		patch, err = translatePatch(patch)
		check(err)

		err = writePatchFile(patch)
		check(err)
	}

	fmt.Printf("Finished in %s", time.Since(start))
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

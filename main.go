package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
)

type engineType int

const (
	engineNone engineType = iota
	engineRPGMVX
	engineWolf
)

var engine = engineNone

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
	file, err := os.Open(filepath.Join(dir, "RPGMKTRANSPATCH"))
	if err != nil {
		file, err = os.Open(filepath.Join(dir, "Patch", "dump", "GameDat.txt"))
		if err != nil {
			return fmt.Errorf("Unable to open RPGMKTRANSPATCH or Patch/dump/GameDat.txt")
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()

		if text == "> RPGMAKER TRANS PATCH V3" {
			fmt.Println("Detected RPG Maker VX Ace Patch")

			engine = engineRPGMVX

			return nil
		} else if text == "> WOLF TRANS PATCH FILE VERSION 1.0" {
			fmt.Println("Detected WOLF RPG Patch")

			engine = engineWolf

			return nil
		}

		return fmt.Errorf("Unsupported patch version")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return err
}

func getDirectoryContents(dir string) []string {
	var fileList []string

	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() || filepath.Ext(path) != ".txt" {
			return nil
		}

		fileList = append(fileList, path)

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return fileList
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

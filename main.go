package main

import (
	"bufio"
	"flag"
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

var (
	engine = engineNone

	// Flags
	lineLength    int
	lineTolerance int
)

func main() {
	start := time.Now()

	args := parseFlags()
	if len(args) < 1 {
		log.Fatal("Program requires patch directory as argument")
	}

	dir := args[0]
	err := checkPatchVersion(dir)
	if err != nil {
		log.Fatal(err)
	}

	fileList := getDirectoryContents(filepath.Join(dir, "Patch"))
	if len(fileList) < 1 {
		log.Fatal("Couldn't find anything to translate")
	}

	if lineLength == -1 {
		if engine == engineWolf {
			lineLength = 54
		} else {
			lineLength = 42
		}
	}

	fmt.Println("Current settings:")
	fmt.Println("- line length:", lineLength)
	fmt.Println("- line length tolerance:", lineTolerance)

	fileCount := len(fileList)

	fmt.Printf("found %d files to translate\n", fileCount)

	jobs, results := createFileWorkers(fileCount)

	go func() {
		for _, file := range fileList {
			jobs <- file
		}
		close(jobs)
	}()

	for err := range results {
		if err != nil {
			log.Error(err)
		}
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

func parseFlags() []string {
	flag.IntVar(&lineLength, "length", -1, "Max line legth")
	flag.IntVar(&lineTolerance, "tolerance", 5, "Max amount of characters allowed to go over the line limit")

	flag.Parse()

	return flag.Args()
}

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestPatchFileParsing(t *testing.T) {
	fileList := getDirectoryContents(filepath.Join("testdata", "Patch"))
	if len(fileList) < 1 {
		t.Error("Couldn't find any files to test")
	}

	for _, inputFile := range fileList {
		patch, err := parsePatchFile(inputFile)
		check(err)

		file, err := ioutil.TempFile(os.TempDir(), "")
		check(err)

		info, err := file.Stat()
		check(err)

		outputFile := filepath.Join(os.TempDir(), info.Name())

		patch.path = outputFile

		err = writePatchFile(patch)
		check(err)

		input, err := ioutil.ReadFile(inputFile)
		check(err)
		output, err := ioutil.ReadFile(outputFile)
		check(err)

		err = os.Remove(patch.path)
		check(err)

		if !bytes.Equal(input, output) {
			t.Errorf("Patch parser couldn't output file equal to input %q", inputFile)
		}
	}
}

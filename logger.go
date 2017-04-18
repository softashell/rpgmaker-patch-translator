package main

import (
	"fmt"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
)

var (
	logMutex = &sync.Mutex{}
)

func logBlockError(err error, args ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	var out string

	out += fmt.Sprintf("Failed to handle block with error %s\nContents:\n", err)

	for _, arg := range args {
		out += fmt.Sprintf("%s\n", spew.Sdump(arg))
	}

	log.Error(out)

	f, err := os.OpenFile("errors.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Error("Unable to open error log file", err)
		return
	}
	defer f.Close()

	if _, err = f.WriteString(out); err != nil {
		log.Error("Unable to write in error log file", err)
	}
}

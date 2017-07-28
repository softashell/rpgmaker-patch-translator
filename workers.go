package main

import (
	"runtime"
	"time"

	"gopkg.in/vbauerster/mpb.v2"
)

func createFileWorkers(fileCount int) (chan string, chan error) {
	workerCount := runtime.NumCPU() / 2

	if workerCount < 1 {
		workerCount = 1
	} else if workerCount > fileCount {
		workerCount = fileCount
	}

	jobs := make(chan string, workerCount)
	results := make(chan error, workerCount)

	p := mpb.New().RefreshRate(100 * time.Millisecond)

	var fileCounter int

	bar := p.AddBar(int64(fileCount)).
		PrependName("All file progress", 25, mpb.DwidthSync|mpb.DextraSpace).
		PrependCounters("%4s/%4s", 0, 10, mpb.DwidthSync|mpb.DextraSpace)

	// Start workers
	for w := 1; w <= workerCount; w++ {
		go func(jobs <-chan string, results chan<- error) {
			for j := range jobs {
				fileCounter++

				results <- processFile(p, fileCounter, fileCount, j)

				bar.Incr(1)
			}
			// not going to get any more jobs, remove worker and close result channel if it was the last worker
			workerCount--
			if workerCount < 1 {
				close(results)
				p.Stop()
			}
		}(jobs, results)
	}

	return jobs, results
}

func processFile(p *mpb.Progress, fileNum, fileCount int, file string) error {
	patch, err := parsePatchFile(file)
	if err != nil {
		return err
	}

	patch, err = translatePatch(p, fileNum, fileCount, patch)
	if err != nil {
		return err
	}

	err = writePatchFile(patch)
	if err != nil {
		return err
	}

	return nil
}

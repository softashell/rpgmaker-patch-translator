package main

import (
	"sync"
	"time"

	"gopkg.in/vbauerster/mpb.v2"
)

type blockWork struct {
	id    int // Only needed to preserve order in patch file
	block patchBlock
}

func createFileWorkers(fileCount int) (chan string, chan error) {
	workerCount := cFileThreads

	if workerCount < 1 {
		workerCount = 1
	} else if workerCount > fileCount {
		workerCount = fileCount
	}

	jobs := make(chan string, workerCount)
	results := make(chan error, workerCount)

	p := mpb.New().RefreshRate(100 * time.Millisecond)

	bar := p.AddBar(int64(fileCount)).
		PrependName("Overall progress", 25, mpb.DwidthSync|mpb.DextraSpace).
		PrependCounters("%4s/%4s", 0, 10, mpb.DwidthSync|mpb.DextraSpace)

	lock := sync.Mutex{}

	// Start workers
	for w := 1; w <= workerCount; w++ {
		go func(jobs <-chan string, results chan<- error) {
			for j := range jobs {
				results <- processFile(p, j)
				bar.Incr(1)
			}

			lock.Lock()
			defer lock.Unlock()

			workerCount--
			if workerCount < 1 {
				close(results)
				p.Stop()
			}
		}(jobs, results)
	}

	return jobs, results
}

func createBlockWorkers(fileCount int) (chan blockWork, chan blockWork) {
	workerCount := cBlockThreads

	if workerCount < 1 {
		workerCount = 1
	} else if workerCount > fileCount {
		workerCount = fileCount
	}

	lock := sync.Mutex{}

	jobs := make(chan blockWork, workerCount)
	results := make(chan blockWork, workerCount)

	for w := 1; w <= workerCount; w++ {
		go func(jobs <-chan blockWork, results chan<- blockWork) {
			for j := range jobs {
				j.block = parseBlock(j.block)
				results <- j
			}

			lock.Lock()
			defer lock.Unlock()

			workerCount--
			if workerCount < 1 {
				close(results)
			}
		}(jobs, results)
	}

	return jobs, results
}

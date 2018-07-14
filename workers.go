package main

import (
	"sync"
	"time"

	"gitgud.io/softashell/rpgmaker-patch-translator/block"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

type blockWork struct {
	id    int // Only needed to preserve order in patch file
	block block.PatchBlock
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

	p := mpb.New(
		mpb.WithRefreshRate(100 * time.Millisecond),
	)

	bar := p.AddBar(int64(fileCount),
		mpb.PrependDecorators(
			decor.Name("Overall progress", decor.WC{W: 25, C: decor.DSyncSpace}),
			decor.CountersNoUnit("%d / %d", decor.WC{C: decor.DSyncSpace}),
		),
	)

	lock := sync.Mutex{}

	// Start workers
	for w := 1; w <= workerCount; w++ {
		go func(jobs <-chan string, results chan<- error) {
			for j := range jobs {
				results <- processFile(p, j)
				bar.Increment()
			}

			lock.Lock()
			defer lock.Unlock()

			workerCount--
			if workerCount < 1 {
				close(results)
				p.Wait()
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

	jobs := make(chan blockWork, workerCount)
	results := make(chan blockWork, workerCount)

	lock := sync.Mutex{}

	for w := 1; w <= workerCount; w++ {
		go func(jobs <-chan blockWork, results chan<- blockWork) {
			for j := range jobs {
				j.block = block.ParseBlock(j.block)
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

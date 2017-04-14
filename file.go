package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/vbauerster/mpb.v2"
)

type patchFile struct {
	path    string
	version string
	blocks  []patchBlock
}

func writePatchFile(patch patchFile) error {
	log.Debugf("Writing %s", patch.path)

	err := os.Remove(patch.path)
	check(err)

	f, err := os.Create(patch.path)
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)

	_, err = w.WriteString(fmt.Sprintf("> %s\n", patch.version))
	check(err)

	for _, block := range patch.blocks {
		_, err = w.WriteString("> BEGIN STRING\n")
		check(err)

		_, err = w.WriteString(block.original)
		check(err)

		for _, context := range block.contexts {
			context = fmt.Sprintf("> CONTEXT: %s", context)

			if !block.translated {
				context += " < UNTRANSLATED\n"
			} else {
				context += "\n"
			}
			_, err = w.WriteString(context)
			check(err)
		}

		var trans string

		if block.translated {
			if block.touched && shouldBreakLines(block.contexts) {
				trans = breakLines(block.translation)
				if !strings.HasSuffix(trans, "\n") {
					trans += "\n"
				}
			} else {
				trans = block.translation
			}
		} else {
			trans = "\n"
		}
		_, err = w.WriteString(trans)
		check(err)

		_, err = w.WriteString("> END STRING\n\n")
		check(err)
	}

	err = w.Flush()
	check(err)

	log.Debugf("Done writing %s", patch.path)

	return nil
}

func parsePatchFile(file string) (patchFile, error) {
	log.Debugf("Parsing %q", path.Base(file))

	f, err := os.Open(file)
	defer f.Close()

	patch := patchFile{path: file}

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)

	original := false
	translation := false

	var orig string
	var trans string
	var contexts []string

	var block patchBlock

	for s.Scan() {
		l := s.Text()

		if strings.HasPrefix(l, "> ") {
			l = l[2:]

			switch {
			case strings.HasPrefix(l, "RPGMAKER TRANS PATCH FILE VERSION"):
				patch.version = l
			case strings.HasPrefix(l, "BEGIN STRING"):
				original = true
			case strings.HasPrefix(l, "CONTEXT: "):
				if len(l) > len("CONTEXT: ")+1 {
					start := len("CONTEXT: ")
					end := strings.Index(l, " < UNTRANSLATED")
					if end == -1 {
						contexts = append(contexts, l[start:])
					} else {
						contexts = append(contexts, l[start:end])
					}
				} else {
					log.Warn("Empty context?", l)
				}

				original = false
				translation = true
			case strings.HasPrefix(l, "END STRING"):
				translation = false

				block.original = orig
				block.contexts = contexts

				if len(strings.TrimRight(trans, "\n")) < 1 {
					block.translation = ""
				} else {
					block.translated = true
					block.translation = trans
				}

				//log.Info(spew.Sdump(block))

				patch.blocks = append(patch.blocks, block)

				orig = ""
				trans = ""
				contexts = nil
			default:
				log.Warn("Unknown input:", l)
			}

			continue
		}

		if !strings.HasSuffix(l, "\n") && (original || translation) {
			l += "\n"
		}

		if original {
			orig += l
		} else if translation {
			trans += l
		}
	}

	return patch, err
}

func translatePatch(patch patchFile) (patchFile, error) {
	var err error

	/* TODO: Add option to skip unsafe scripts (everything but vocab)
	if strings.HasSuffix(patch.path, "Scripts.txt") {
		fmt.Println("Skipped")

		return patch, err
	}
	*/

	// Only needed to preserve order in patch file
	type blockWork struct {
		id    int
		block patchBlock
	}

	jobs := make(chan blockWork, runtime.NumCPU()*2)
	results := make(chan blockWork, runtime.NumCPU()*2)

	// Start workers
	for w := 1; w <= runtime.NumCPU(); w++ {
		go func(jobs <-chan blockWork, results chan<- blockWork) {
			for j := range jobs {
				j.block = parseBlock(j.block)
				results <- j
			}
		}(jobs, results)
	}

	p := mpb.New().
		RefreshRate(500 * time.Millisecond)

	defer p.Stop()

	count := len(patch.blocks)

	bar := p.AddBar(int64(count)).
		PrependCounters("%4s/%4s", 0, 10, mpb.DwidthSync|mpb.DextraSpace).
		AppendETA(2, 0)

	// Add blocks in background to job queue
	go func() {
		for i, block := range patch.blocks {
			//patch.blocks[i] = parseBlock(block)
			w := blockWork{i, block}
			jobs <- w
		}
		close(jobs)
	}()

	// Start reading results, will block if there are none
	for a := count; a > 0; a-- {

		j := <-results

		patch.blocks[j.id] = j.block

		bar.Incr(1)
	}

	return patch, err
}

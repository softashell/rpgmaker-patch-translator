package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dimchansky/utfbom"
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

		for _, t := range block.translations {
			for _, context := range t.contexts {
				context = fmt.Sprintf("> CONTEXT: %s", context)

				if !t.translated {
					context += " < UNTRANSLATED\n"
				} else {
					context += "\n"
				}
				_, err = w.WriteString(context)
				check(err)
			}

			var trans string

			if t.translated {
				if t.touched && shouldBreakLines(t.contexts) {
					trans = breakLines(t.text)
					if !strings.HasSuffix(trans, "\n") {
						trans += "\n"
					}
				} else {
					trans = t.text
				}
			} else {
				trans = "\n"
			}

			_, err = w.WriteString(trans)
			check(err)
		}

		_, err = w.WriteString("> END STRING\n\n")
		check(err)
	}

	err = w.Flush()
	check(err)

	log.Debugf("Done writing %s", patch.path)

	return nil
}

func parsePatchFile(file string) (patchFile, error) {
	log.Debugf("Parsing %q", filepath.Base(file))

	patch := patchFile{path: file}

	f, err := os.Open(file)
	if err != nil {
		return patch, err
	}
	defer f.Close()

	s := bufio.NewScanner(utfbom.SkipOnly(f))
	s.Split(bufio.ScanLines)

	original := false
	translation := false

	var orig string
	var trans string
	var contexts []string

	var block patchBlock
	var translations []translationBlock

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
				if translation && len(trans) > 0 {
					var translated bool

					if len(strings.TrimRight(trans, "\n")) < 1 {
						trans = ""
						translated = false
					} else {
						translated = true
					}

					translations = append(translations, translationBlock{
						text:       trans,
						contexts:   contexts,
						translated: translated,
					})

					trans = ""
					contexts = nil
					translation = false
				} else {
					original = false
					translation = true
				}

				if len(l) > len("CONTEXT: ")+1 {
					start := len("CONTEXT: ")
					end := strings.Index(l, " < UNTRANSLATED")
					if end == -1 {
						contexts = append(contexts, l[start:])
					} else {
						contexts = append(contexts, l[start:end])
					}
				}
			case strings.HasPrefix(l, "END STRING"):
				translation = false

				if len(trans) > 0 {
					var translated bool

					if len(strings.TrimRight(trans, "\n")) < 1 {
						trans = ""
						translated = false
					} else {
						translated = true
					}

					translations = append(translations, translationBlock{
						text:       trans,
						contexts:   contexts,
						translated: translated,
					})
				}

				if len(translations) == 0 {
					log.Error("No contexts in block")
				}

				block.original = orig
				block.translations = translations

				patch.blocks = append(patch.blocks, block)

				orig = ""
				trans = ""

				contexts = nil
				translations = nil
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

	if len(patch.version) < 3 {
		err = fmt.Errorf("No patch version found in %q", file)
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

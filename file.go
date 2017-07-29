package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitgud.io/softashell/rpgmaker-patch-translator/text"

	log "github.com/Sirupsen/logrus"
	"github.com/dimchansky/utfbom"
	"github.com/pkg/errors"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
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

		_, err = w.WriteString(text.Escape(block.original))
		check(err)

		for _, t := range block.translations {
			for _, context := range t.contexts {
				context = fmt.Sprintf("> CONTEXT%s", context)

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
				text := text.Escape(t.text)

				if t.touched && shouldBreakLines(t.contexts) {
					trans = breakLines(text)
					if !strings.HasSuffix(trans, "\n") {
						trans += "\n"
					}
				} else {
					trans = text
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
		return patch, errors.Wrapf(err, "failed to open patch file: %q", file)
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
			case strings.HasPrefix(l, "RPGMAKER TRANS PATCH FILE VERSION") || strings.HasPrefix(l, "WOLF TRANS PATCH FILE VERSION 1.0"):
				patch.version = l
			case strings.HasPrefix(l, "BEGIN STRING"):
				original = true
			case strings.HasPrefix(l, "CONTEXT"):
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

				if len(l) > len("CONTEXT")+1 {
					start := len("CONTEXT")
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
						text:       text.Unescape(trans),
						contexts:   contexts,
						translated: translated,
					})
				} else if len(contexts) > 0 {
					translations = append(translations, translationBlock{
						text:       text.Unescape(trans),
						contexts:   contexts,
						translated: false,
					})
				} else if len(translations) == 0 {
					log.Errorf("No contexts found for block with original text:\n%q", orig)
				}

				block.original = text.Unescape(orig)
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

	if err := s.Err(); err != nil {
		return patch, errors.Wrapf(err, "error while scanning patch file: %q", file)
	}

	if len(patch.version) < 3 {
		err = fmt.Errorf("No patch version found in %q", file)
	}

	return patch, err
}

func translatePatch(p *mpb.Progress, patch patchFile) (patchFile, error) {
	blockCount := len(patch.blocks)

	jobs, results := createBlockWorkers(blockCount)

	bar := p.AddBar(int64(blockCount),
		mpb.PrependDecorators(
			decor.StaticName(filepath.Base(patch.path), 25, decor.DwidthSync|decor.DextraSpace),
			decor.Counters("%4s/%4s", 0, 10, decor.DwidthSync|decor.DextraSpace),
		))

	// Add blocks in background to job queue
	go func() {
		for i, block := range patch.blocks {
			w := blockWork{i, block}
			jobs <- w
		}
		close(jobs)
	}()

	// Start reading results, will block if there are none
	for j := range results {
		patch.blocks[j.id] = j.block
		bar.Incr(1)
	}

	p.RemoveBar(bar)

	return patch, nil
}

func processFile(p *mpb.Progress, file string) error {
	patch, err := parsePatchFile(file)
	if err != nil {
		return err
	}

	patch, err = translatePatch(p, patch)
	if err != nil {
		return err
	}

	err = writePatchFile(patch)
	if err != nil {
		return err
	}

	return nil
}

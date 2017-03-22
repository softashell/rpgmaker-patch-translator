package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type patchFile struct {
	path    string
	version string
	blocks  []patchBlock
}

func writePatchFile(patch patchFile) error {
	log.Infof("Writing %s", patch.path)

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

		if block.translated {
			trans := breakLines(block.translation)

			_, err = w.WriteString(trans)
			check(err)
		}

		_, err = w.WriteString("> END STRING\n\n")
		check(err)
	}

	err = w.Flush()
	check(err)

	log.Infof("Done writing %s", patch.path)

	return nil
}

func parsePatchFile(path string) (patchFile, error) {
	log.Info("Parsing", path)

	f, err := os.Open(path)
	defer f.Close()

	file := patchFile{path: path}

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
				file.version = l
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

				file.blocks = append(file.blocks, block)

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

	return file, err
}

func translatePatch(patch patchFile) (patchFile, error) {
	var err error

	if strings.HasSuffix(patch.path, "Scripts.txt") {
		return patch, err
	}

	for i, block := range patch.blocks {
		patch.blocks[i] = parseBlock(block)
	}

	return patch, err
}

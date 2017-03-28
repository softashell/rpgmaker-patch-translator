package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
)

type patchBlock struct {
	original    string
	contexts    []string
	translation string
	translated  bool
}

func parseBlock(block patchBlock) patchBlock {
	if shouldTranslate(block) {
		items, err := parseText(block.original)
		if err != nil {
			log.Errorf("%s\ncontexts:\n%v\ntext: %q", err, block.contexts, block.original)
			log.Error(spew.Sdump(items))

			// TODO: Avoid panic here and just log the offending block to file along with debug info
			panic(err)
		}

		block.translation = translateItems(items)
		block.translated = true

		log.Debugf("'%s' => '%s'\n", block.original, block.translation)
	}

	return block
}

func shouldTranslate(block patchBlock) bool {
	if block.translated {
		log.Debug("Skipping translated block")
		return false
	}

	for _, c := range block.contexts {
		log.Debugf("%q", c)
		if strings.HasSuffix(c, "_se/name/") ||
			strings.HasSuffix(c, "/bgm/name/") ||
			strings.HasSuffix(c, "_me/name/") ||
			strings.Contains(c, "/InlineScript/") {
			return false
		}
	}

	return true
}

func translateItems(items []item) string {
	var out string

	for _, item := range items {
		log.Debugf("%14s: %q", item.typ, item.val)

		// Translate
		if item.typ == itemText {
			translation, err := translateString(item.val)
			check(err)

			if strings.HasPrefix(item.val, " ") && !strings.HasPrefix(translation, " ") {
				out += " "
			}

			out += translation

			if strings.HasSuffix(item.val, " ") && !strings.HasSuffix(translation, " ") {
				out += " "
			}
		} else if item.typ == itemEOF {
			continue
		} else if item.typ != itemError {
			// Add space before '(' since translation might make it get parsed as function
			if item.typ == itemRawString && item.val == "(" &&
				(!strings.HasSuffix(out, " ") && !strings.HasSuffix(out, "　")) {
				out += " "
			}

			// Add raw
			out += item.val
		}
	}

	return out
}

func parseText(text string) ([]item, error) {
	l := lex(text)

	var items []item
	var err error

	for {
		item := l.nextItem()
		items = append(items, item)

		if item.typ == itemError {
			err = fmt.Errorf("Failed to parse text")
		}

		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}

	return items, err
}

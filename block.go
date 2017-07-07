package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
)

type patchBlock struct {
	original     string
	translations []translationBlock
}

type translationBlock struct {
	contexts   []string
	text       string
	touched    bool
	translated bool
}

func parseBlock(block patchBlock) patchBlock {
	if !shouldTranslateText(block.original) {
		return block
	}

	var err error
	var items []item
	var untranslated []string
	var translated, parsed bool

	for i, t := range block.translations {
		if t.translated {
			continue // Block is already translated
		}

		good, bad := getTranslatableContexts(t, block.original)
		untranslated = append(untranslated, bad...)

		if len(good) > 0 {
			if !parsed {
				items, err = parseText(block.original)

				if err != nil {
					logBlockError(err, block, items)

					return block
				}

				parsed = true
			}

			t.text = translateItems(items)
			t.translated = true
			t.touched = true
			t.contexts = good

			block.translations[i] = t

			log.Debugf("'%s' => '%s'\n", block.original, t.text)

			translated = true
		}
	}

	if translated && len(untranslated) > 0 {
		block.translations = append(block.translations, translationBlock{
			text:       "",
			contexts:   untranslated,
			translated: false,
		})

		log.Warnf("Mixed block\n %s", spew.Sdump(block))
	}

	return block
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
			break
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

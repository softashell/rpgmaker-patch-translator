package main

import (
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
	if shouldTranslate(block) && len(block.translation) < 2 {
		items := parseText(block.original)

		block.translation = translateItems(items)
		block.translated = true

		log.Infof("'%s' => '%s'\n", block.original, block.translation)
	}

	//var input string
	//fmt.Scanln(&input)

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
			strings.HasSuffix(c, "_me/name/") {
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

func parseText(text string) []item {
	l := lex(text)

	var items []item

	error := false

	for {
		item := l.nextItem()
		items = append(items, item)

		if item.typ == itemError {
			error = true
		}

		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}

	if error {
		log.Error(text)
		log.Error(spew.Sdump(items))
		panic("Failed to parse text")
	}

	return items
}

package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
)

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

func translateItems(items []item) string {
	for i := range items {
		if items[i].typ == itemText {
			translation, err := translateString(items[i].val)
			if err != nil {
				log.Fatal("failed to translate item", err)
			}

			if strings.HasPrefix(items[i].val, " ") && !strings.HasPrefix(translation, " ") {
				translation = " " + translation
			}

			if strings.HasSuffix(items[i].val, " ") && !strings.HasSuffix(translation, " ") {
				translation += " "
			}

			items[i].trans = translation
		}
	}

	return assembleItems(items)
}

func assembleItems(items []item) string {
	var out string

	var lastVal string
	var lastType itemType

	for _, item := range items {
		log.Debugf("%14s: %q", item.typ, item.val)

		if item.typ == itemText {
			// Space after raw strings that may contain english and any scripts
			if (lastType == itemRawString && strings.ContainsAny(ignoredCharacters, lastVal)) ||
				(lastType == itemScript || lastType == itemRightDelim || lastType == itemRightParen) {
				if !endsWithWhitespace(out) {
					out += " "
				}
			}

			out += item.trans
		} else if item.typ == itemEOF {
			break
		} else if item.typ != itemError {
			if item.typ == itemRawString {
				if item.val == "(" && !endsWithWhitespace(out) {
					// Add space before '(' since translation might make it get parsed as function
					out += " "
				} else if lastType == itemText && strings.ContainsAny(ignoredCharacters, item.val) {
					// If last item was translated check if we're trying to add something,
					// that might be in english or a number right after it
					if !endsWithWhitespace(out) {
						out += " "
					}
				}
			} else if item.typ == itemScript && lastType == itemText {
				if !endsWithWhitespace(out) {
					out += " "
				}
			}

			// Add raw
			out += item.val
		}

		lastVal = item.val
		lastType = item.typ
	}

	return out
}

package lex

import (
	"fmt"
	"strings"

	"golang.org/x/text/width"

	"gitgud.io/softashell/rpgmaker-patch-translator/text"
	"gitgud.io/softashell/rpgmaker-patch-translator/translate"
	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

func ParseText(text string) ([]Item, error) {
	l := lex(text)

	var items []Item
	var err error

	for {
		item := l.nextItem()
		items = append(items, item)

		if item.Typ == ItemError {
			err = fmt.Errorf("Failed to parse text: %s", item.Val)
		}

		if item.Typ == ItemEOF || item.Typ == ItemError {
			break
		}
	}

	return items, err
}

func getOnlyText(text string) string {
	items, err := ParseText(text)
	if err != nil {
		log.Errorf("%s\ntext: %q", err, text)
		log.Error(spew.Sdump(items))

		panic(err)
	}

	var out string

	for _, item := range items {
		switch item.Typ {
		case ItemText, ItemRawString, ItemNumber:
			out += item.Val
		}
	}

	return out
}

func assembleItems(items []Item) string {
	var out string

	var lastVal string
	var lastType itemType

	for _, item := range items {
		log.Debugf("%14s: %q", item.Typ, item.Val)

		if item.Typ == ItemText {
			// Space after raw strings that may contain english and any scripts
			if (lastType == ItemRawString && strings.ContainsAny(ignoredCharacters, lastVal)) ||
				(lastType == ItemScript || lastType == ItemRightDelim || lastType == ItemRightParen) {
				if !text.EndsWithWhitespace(out) && !text.StartsWithWhitespace(item.Val) {
					out += " "
				}
			} else if lastType == ItemNumber {
				if !text.EndsWithWhitespace(out) {
					out += " "
				}
			}

			out += item.Val
		} else if item.Typ == ItemEOF {
			break
		} else if item.Typ != ItemError {
			if item.Typ == ItemRawString {
				if item.Val == "(" && !text.EndsWithWhitespace(out) {
					// Add space before '(' since translation might make it get parsed as function
					out += " "
				} else if lastType == ItemText && strings.ContainsAny(ignoredCharacters, item.Val) {
					// If last item was translated check if we're trying to add something,
					// that might be in english or a number right after it
					if !text.EndsWithWhitespace(out) && !text.StartsWithWhitespace(item.Val) {
						out += " "
					}
				}
			} else if item.Typ == ItemScript && (lastType == ItemText || lastType == ItemNumber) {
				if !text.EndsWithWhitespace(out) {
					out += " "
				}
			} else if item.Typ == ItemNumber && lastType == ItemText {
				if !text.EndsWithWhitespace(out) {
					out += " "
				}
			}

			// Add raw
			out += item.Val
		}

		lastType = item.Typ
		lastVal = item.Val
	}

	return out
}

func TranslateItems(items []Item) (string, error) {
	for i := range items {
		if items[i].Typ == ItemText {
			translation, err := translate.String(items[i].Val)
			if err != nil {
				return "", errors.Wrapf(err, "failed to translate [%s] %q", items[i].Typ, items[i].Val)
			}

			if text.StartsWithWhitespace(items[i].Val) && !text.StartsWithWhitespace(translation) {
				translation = " " + translation
			}

			if text.EndsWithWhitespace(items[i].Val) && !text.EndsWithWhitespace(translation) {
				translation += " "
			}

			items[i].Val = translation
		} else if items[i].Typ == ItemNumber {
			text := items[i].Val
			text = width.Narrow.String(text)

			text, err := translate.String(text)
			if err != nil {
				return "", errors.Wrapf(err, "failed to translate [%s] %q", items[i].Typ, items[i].Val)
			}

			text = strings.ToLower(text)

			items[i].Val = text
		}
	}

	return assembleItems(items), nil
}

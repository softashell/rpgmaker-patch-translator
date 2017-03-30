package main

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
)

func getOnlyText(text string) string {
	items, err := parseText(text)
	if err != nil {
		log.Errorf("%s\ntext: %q", err, text)
		log.Error(spew.Sdump(items))

		panic(err)
	}

	var out string

	for _, item := range items {
		switch item.typ {
		case itemText, itemRawString:
			out += item.val
		}
	}

	return out
}

func breakLines(text string) string {
	var out string

	lines := strings.Split(text, "\n")

	for n, l := range lines {
		out += breakLine(l)

		// Add newline if it's not the last line
		if n+1 < len(lines) {
			out += "\n"

		}
	}

	return out
}

func breakLine(text string) string {
	// Remove any extra trailing new lines
	text = strings.TrimRight(text, "\n")
	if len(text) < 1 {
		return text
	}

	items, err := parseText(text)
	if err != nil {
		log.Errorf("%s\ntext: %q", err, text)
		log.Error(spew.Sdump(items))

		panic(err)
	}

	log.Debug(spew.Sdump(items))

	var out string
	var justText string
	var line string

	for _, item := range items {
		switch item.typ {
		case itemText, itemRawString:
			if len([]rune(justText+item.val)) <= 42 {
				out += line
				out += item.val

				line = ""
				justText += item.val

				break
			}

			log.Debugf("Trying to split %q from %q", item.val, text)

			words := strings.Split(item.val, " ")
			for i := range words {
				if len([]rune(justText+words[i])) <= 42 {
					line += words[i]
					justText += words[i]

					if i+1 < len(words) {
						line += " "
						justText += " "
					}

					continue
				}

				line = strings.TrimRight(line, " ")

				log.Debugf("Split! %q from %q", line, text)

				if !strings.HasSuffix(line, "\n") {
					line += "\n"
				}

				out += line

				line = ""
				justText = ""

			}
		default:
			out += item.val
		}
	}

	if len(line) > 0 {
		log.Debugf("Split! Trailing %q from %q", line, text)

		out += line
	}

	return out
}

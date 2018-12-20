package main

import (
	"strings"

	"gitgud.io/softashell/rpgmaker-patch-translator/lex"
	"gitgud.io/softashell/rpgmaker-patch-translator/text"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

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

func breakLine(input string) string {
	// Remove any extra trailing new lines
	input = strings.TrimRight(input, "\n")
	if len(input) < 1 {
		return input
	}

	items, err := lex.ParseText(input)
	if err != nil {
		log.Errorf("%s\ntext: %q", err, input)
		log.Error(spew.Sdump(items))

		panic(err)
	}

	log.Debug(spew.Sdump(items))

	var out, justText, line string

	leadingWhitespace := text.ExtractLeadingWhitespace(input)

	for _, item := range items {
		switch item.Typ {
		case lex.ItemText, lex.ItemRawString, lex.ItemNumber:
			if len([]rune(justText+item.Val)) <= lineLength {
				out += line
				out += item.Val

				line = ""
				justText += item.Val

				break
			}

			log.Debugf("Trying to split %q from %q", item.Val, input)

			words := strings.Split(item.Val, " ")
			for i := range words {
				log.Debug("word: ", i+1, " / ", len(words), " len:", len([]rune(justText+words[i])))

				if len([]rune(justText+words[i])) <= lineLength {
					log.Debugf("adding %q", words[i])

					line += words[i]
					justText += words[i]

					if i+1 < len(words) {
						line += " "
						justText += " "
					}

					continue
				}

				if i+1 == len(words) && len([]rune(justText+words[i])) <= lineLength+lineTolerance {
					log.Debugf("Word %q was too long to fit! Not adding a new line before because it's short and last item", words[i])
					line += words[i]

					break
				}

				log.Debugf("Word %q was too long to fit! Added a new line before it", words[i])
				line = strings.TrimRight(line, " ")
				line += "\n" + leadingWhitespace

				out += line

				line = words[i]
				justText = words[i]

				if i+1 < len(words) {
					line += " "
					justText += " "
				}
			}
		default:
			line += item.Val
		}
	}

	if len(line) > 0 {
		log.Debugf("Split! Trailing %q from %q", line, input)

		out += line
	}

	return out
}

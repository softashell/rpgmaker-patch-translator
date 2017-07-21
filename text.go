package main

import (
	"regexp"
	"strings"

	"golang.org/x/text/width"

	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
)

var (
	moonRegex         = regexp.MustCompile(`(\p{Hiragana}|\p{Katakana}|\p{Han})`)
	ignoredExtensions = []string{
		".png",
		".jpg",
		".jpeg",
		".gif",
		".bmp",
		".ogg",
		".mp3",
		".wav",
		".mid",
		".midi",
		".txt",
		".csv",
	}
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
		case itemText, itemRawString, itemNumber:
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

	var out, justText, line string

	for _, item := range items {
		switch item.typ {
		case itemText, itemRawString, itemNumber:
			if len([]rune(justText+item.val)) <= lineLength {
				out += line
				out += item.val

				line = ""
				justText += item.val

				break
			}

			log.Debugf("Trying to split %q from %q", item.val, text)

			words := strings.Split(item.val, " ")
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
				line += "\n"

				out += line

				line = words[i]
				justText = words[i]

				if i+1 < len(words) {
					line += " "
					justText += " "
				}
			}
		default:
			line += item.val
		}
	}

	if len(line) > 0 {
		log.Debugf("Split! Trailing %q from %q", line, text)

		out += line
	}

	return out
}

func shouldTranslateText(text string) bool {
	text = strings.TrimSpace(text)

	if len(text) < 1 {
		return false
	}

	//Regex
	if strings.HasPrefix(text, "/") && strings.HasSuffix(text, "/") {
		return false
	}

	text = width.Narrow.String(text)
	text = strings.ToLower(text)

	for _, e := range ignoredExtensions {
		if strings.Contains(text, e) {
			return false
		}
	}

	return isJapanese(text)
}

func isJapanese(text string) bool {
	matches := moonRegex.FindAllString(text, 1)

	return len(matches) >= 1
}

func replaceRegex(text string, expr string, repl string) string {
	regex := regexp.MustCompile(expr)

	return regex.ReplaceAllString(text, repl)
}

func unescapeText(text string) string {
	text = strings.Replace(text, `\\`, `\`, -1)

	text = strings.Replace(text, `\>`, `>`, -1)
	text = strings.Replace(text, `\#`, `#`, -1)

	return text
}

func escapeText(text string) string {
	text = strings.Replace(text, `\`, `\\`, -1)

	text = strings.Replace(text, `>`, `\>`, -1)
	text = strings.Replace(text, `#`, `\#`, -1)

	return text
}

func startsWithWhitespace(text string) bool {
	if strings.HasPrefix(text, " ") || strings.HasPrefix(text, "　") {
		return true
	}

	return false
}

func endsWithWhitespace(text string) bool {
	if strings.HasSuffix(text, " ") || strings.HasSuffix(text, "　") {
		return true
	}

	return false
}

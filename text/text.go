package text

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/width"
)

var (
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

func ShouldTranslate(text string) bool {
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
	for _, r := range text {
		if unicode.Is(unicode.Hiragana, r) ||
			unicode.Is(unicode.Katakana, r) ||
			unicode.Is(unicode.Han, r) {
			return true
		}
	}

	return false
}

func ReplaceRegex(text string, expr string, repl string) string {
	regex := regexp.MustCompile(expr)

	return regex.ReplaceAllString(text, repl)
}

func Unescape(text string) string {
	text = strings.Replace(text, `\\`, `\`, -1)

	text = strings.Replace(text, `\>`, `>`, -1)
	text = strings.Replace(text, `\#`, `#`, -1)

	return text
}

func Escape(text string) string {
	text = strings.Replace(text, `\`, `\\`, -1)

	text = strings.Replace(text, `>`, `\>`, -1)
	text = strings.Replace(text, `#`, `\#`, -1)

	return text
}

func StartsWithWhitespace(text string) bool {
	r, _ := utf8.DecodeRuneInString(text)

	if unicode.IsSpace(r) {
		return true
	}

	return false
}

func EndsWithWhitespace(text string) bool {
	r, _ := utf8.DecodeLastRuneInString(text)

	if unicode.IsSpace(r) {
		return true
	}

	return false
}

func ExtractLeadingWhitespace(s string) string {
	var ws string

	for _, r := range s {
		if unicode.IsSpace(r) {
			ws += string(r)
			continue
		}

		break
	}

	return ws
}

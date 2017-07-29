package text

import (
	"regexp"
	"strings"

	"golang.org/x/text/width"
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
	matches := moonRegex.FindAllString(text, 1)

	return len(matches) >= 1
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
	if strings.HasPrefix(text, " ") || strings.HasPrefix(text, "　") {
		return true
	}

	return false
}

func EndsWithWhitespace(text string) bool {
	if strings.HasSuffix(text, " ") || strings.HasSuffix(text, "　") {
		return true
	}

	return false
}

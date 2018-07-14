package translate

import (
	"net/http"
	"strings"
	"unicode"

	"gitgud.io/softashell/rpgmaker-patch-translator/text"
	log "github.com/Sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
)

type translateRequest struct {
	Text string `json:"text"`
	From string `json:"from"`
	To   string `json:"to"`
}

type translateResponse struct {
	Text            string `json:"text"`
	From            string `json:"from"`
	To              string `json:"to"`
	TranslationText string `json:"translationText"`
}

var httpTransport *http.Transport

func Init() {
	httpTransport = &http.Transport{}
}

func String(str string) (string, error) {
	if !text.ShouldTranslate(str) {
		return str, nil
	}

	if str == "っ" {
		return "", nil
	}

	var response translateResponse

	request := translateRequest{
		From: "ja",
		To:   "en",
		Text: str,
	}

	gr := gorequest.New()
	gr.Transport = httpTransport

	_, _, errs := gr.Post("http://127.0.0.1:3000/api/translate").
		Type("json").SendStruct(&request).EndStruct(&response)
	for _, err := range errs {
		return "", errors.Wrap(err, "http request failed")
	}

	out := response.TranslationText

	if len(out) < 1 {
		log.Warnf("Translator returned empty string, replacing with original text %q", str)
		out = str
	} else {
		out = cleanTranslation(out)
	}

	return out, nil
}

func cleanTranslation(str string) string {
	// Removes any rune that isn't printable or a space
	isValid := func(r rune) rune {
		if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			return -1
		}

		return r
	}

	str = strings.Map(isValid, str)

	str = strings.Replace(str, "\\u0026", "＆", -1)
	str = strings.Replace(str, "\\u003c", "<", -1)
	str = strings.Replace(str, "\\u003e", ">", -1)

	if strings.Contains(str, "\\u0") {
		log.Warnf("Found unexpected escaped character in translation %s", str)
	}

	str = strings.Replace(str, "\\", "", -1)

	// Repeated whitespace
	str = text.ReplaceRegex(str, `\s{2,}`, " ")

	// ー ー ー ー
	str = text.ReplaceRegex(str, `\s+((\s+)?[-―ー]){2,}`, " ー")

	// · · · ·
	str = text.ReplaceRegex(str, `(\s+)?((\s+)?[·]+){3,}`, " ···")

	str = text.ReplaceRegex(str, `((\s+)?っ)+`, "")

	return str
}

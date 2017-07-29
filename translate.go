package main

import (
	"strings"
	"unicode"

	log "github.com/Sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
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

func translateString(text string) (string, error) {
	if !shouldTranslateText(text) {
		return text, nil
	}

	if text == "っ" {
		return "", nil
	}

	var response translateResponse

	request := translateRequest{
		From: "ja",
		To:   "en",
		Text: text,
	}

	_, _, errs := gorequest.New().Post("http://127.0.0.1:3000/api/translate").
		Type("json").SendStruct(&request).EndStruct(&response)
	for _, err := range errs {
		return "", err
	}

	out := response.TranslationText

	if len(out) < 1 {
		log.Warnf("Translator returned empty string, replacing with original text %q", text)
		out = text
	} else {
		out = cleanTranslation(out)
	}

	return out, nil
}

func cleanTranslation(text string) string {
	// Removes any rune that isn't printable or a space
	isValid := func(r rune) rune {
		if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			return -1
		}

		return r
	}

	text = strings.Map(isValid, text)

	text = strings.Replace(text, "\\u0026", "＆", -1)
	text = strings.Replace(text, "\\u003c", "<", -1)
	text = strings.Replace(text, "\\u003e", ">", -1)

	if strings.Contains(text, "\\u0") {
		log.Warnf("Found unexpected escaped character in translation %s", text)
	}

	text = strings.Replace(text, "\\", "", -1)

	// Repeated whitespace
	text = replaceRegex(text, `\s{2,}`, " ")

	// ー ー ー ー
	text = replaceRegex(text, `\s+((\s+)?[-―ー]){2,}`, " ー")

	// · · · ·
	text = replaceRegex(text, `(\s+)?((\s+)?[·]+){3,}`, " ···")

	text = replaceRegex(text, `((\s+)?っ)+`, "")

	return text
}

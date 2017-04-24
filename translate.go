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

	var response translateResponse

	request := translateRequest{
		From: "ja",
		To:   "en",
		Text: text,
	}

	resp, reply, errs := gorequest.New().Post("http://127.0.0.1:3000/api/translate").
		Type("json").SendStruct(&request).EndStruct(&response)
	for _, err := range errs {
		log.WithFields(log.Fields{
			"response": resp,
			"reply":    reply,
		}).Error(err)

		return "", err
	}

	out := response.TranslationText
	out = cleanTranslation(out)

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

	// Repeated whitespace
	text = replaceRegex(text, `\s{2,}`, " ")

	// ー ー ー ー
	text = replaceRegex(text, `\s+([-―ー](\s+)?){2,}`, "―")

	return text
}

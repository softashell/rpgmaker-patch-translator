package main

import (
	"regexp"
	"strings"

	"golang.org/x/text/width"

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

	return out, nil
}

func shouldTranslateText(text string) bool {
	text = strings.TrimSpace(text)

	if len(text) < 1 {
		return false
	}

	text = width.Narrow.String(text)

	if !isJapanese(text) {
		return false
	}

	return true
}

func isJapanese(text string) bool {
	regex := regexp.MustCompile(`(\p{Hiragana}|\p{Katakana}|\p{Han})`)
	matches := regex.FindAllString(text, 1)

	return len(matches) >= 1
}

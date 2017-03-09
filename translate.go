package main

import (
	"strings"

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
	if len(strings.TrimSpace(text)) < 1 {
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

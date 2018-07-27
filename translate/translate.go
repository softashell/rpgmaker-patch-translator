package translate

import (
	"net/http"
	"net/rpc"
	"strings"
	"unicode"

	"gitgud.io/softashell/rpgmaker-patch-translator/text"
	"github.com/Jeffail/tunny"
	log "github.com/Sirupsen/logrus"
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
var pool *tunny.Pool

func Init() {
	httpTransport = &http.Transport{}

	workerCount := 64

	pool = tunny.New(workerCount, func() tunny.Worker {
		return newComfyWorker()
	})
}

func comfyTranslate(client *rpc.Client, req translateRequest) translateResponse {
	var reply translateResponse

	err := client.Call("Comfy.Translate", req, &reply)
	if err != nil {
		log.Fatal("translation service error:", err)
	}

	return reply
}

func String(str string) (string, error) {
	if !text.ShouldTranslate(str) {
		return str, nil
	}

	if str == "っ" {
		return "", nil
	}

	request := translateRequest{
		From: "ja",
		To:   "en",
		Text: str,
	}

	response := pool.Process(request).(translateResponse)

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

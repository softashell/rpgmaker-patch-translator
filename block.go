package main

import (
	"gitgud.io/softashell/rpgmaker-patch-translator/lex"
	"gitgud.io/softashell/rpgmaker-patch-translator/text"
	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
)

type patchBlock struct {
	original     string
	translations []translationBlock
}

type translationBlock struct {
	contexts   []string
	text       string
	touched    bool
	translated bool
}

func parseBlock(block patchBlock) patchBlock {
	if !text.ShouldTranslate(block.original) {
		return block
	}

	var err error
	var items []lex.Item
	var untranslated []string
	var translated, parsed bool

	for i, t := range block.translations {
		if t.translated {
			continue // Block is already translated
		}

		good, bad := getTranslatableContexts(t, block.original)
		untranslated = append(untranslated, bad...)

		if len(good) < 1 {
			continue
		}

		if !parsed {
			items, err = lex.ParseText(block.original)
			if err != nil {
				logBlockError(err, block, items)

				return block
			}

			parsed = true
		}

		t.text, err = lex.TranslateItems(items)
		if err != nil {
			// This should only fail if translation service is down
			log.Fatalf("failed to translate items: %v", err)
		}

		t.translated = true
		t.touched = true
		t.contexts = good

		block.translations[i] = t

		log.Debugf("'%s' => '%s'\n", block.original, t.text)

		translated = true
	}

	if translated && len(untranslated) > 0 {
		block.translations = append(block.translations, translationBlock{
			text:       "",
			contexts:   untranslated,
			translated: false,
		})

		log.Infof("Mixed block\n %s", spew.Sdump(block))
	}

	return block
}

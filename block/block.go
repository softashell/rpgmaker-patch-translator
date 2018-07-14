package block

import (
	"gitgud.io/softashell/rpgmaker-patch-translator/lex"
	"gitgud.io/softashell/rpgmaker-patch-translator/statictl"
	"gitgud.io/softashell/rpgmaker-patch-translator/text"
	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

type PatchBlock struct {
	Original     string
	Translations []TranslationBlock
}

type TranslationBlock struct {
	Contexts   []string
	Text       string
	Touched    bool
	Translated bool
}

var stl *statictl.Db

func Init() {
	// Set the nasty global variables

	if stl == nil {
		stl = statictl.New()
	}
}

func ParseBlock(block PatchBlock) PatchBlock {
	if !text.ShouldTranslate(block.Original) {
		return block
	}

	var err error
	var items []lex.Item
	var untranslated []string
	var translated, parsed bool

	for i, t := range block.Translations {
		if t.Translated {
			continue // Block is already translated
		}

		// Attempt to get static translation for block
		t, err = TranslateBlockStatic(t, block.Original)

		// Fallback to lexing and translating chunks with comfy-translator
		if err != nil {
			good, bad := getTranslatableContexts(t, block.Original)
			untranslated = append(untranslated, bad...)

			if len(good) < 1 {
				continue
			}

			if !parsed {
				items, err = lex.ParseText(block.Original)
				if err != nil {
					//logBlockError(err, block, items)

					return block
				}

				parsed = true
			}

			t.Text, err = lex.TranslateItems(items)
			if err != nil {
				// This should only fail if translation service is down
				log.Fatalf("failed to translate items: %v", err)
			}

			t.Contexts = good
			t.Translated = true
			t.Touched = true
		}

		block.Translations[i] = t

		log.Debugf("'%s' => '%s'\n", block.Original, t.Text)

		translated = true
	}

	if translated && len(untranslated) > 0 {
		block.Translations = append(block.Translations, TranslationBlock{
			Text:       "",
			Contexts:   untranslated,
			Translated: false,
		})

		log.Infof("Mixed block\n %s", spew.Sdump(block))
	}

	return block
}

func TranslateBlockStatic(b TranslationBlock, originalText string) (TranslationBlock, error) {
	// TODO: Match every context against translation and create new blocks if they differ

	tlType := GetTranslationType(b.Contexts)

	text, err := stl.GetTranslation(originalText, tlType)
	if err != nil {
		return b, errors.Wrapf(err, "can't get translation for %q %q", originalText, tlType)
	}

	b.Text = text
	b.Translated = true
	b.Touched = true

	return b, nil
}

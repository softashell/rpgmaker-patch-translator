package block

import (
	"gitgud.io/softashell/rpgmaker-patch-translator/lex"
	"gitgud.io/softashell/rpgmaker-patch-translator/statictl"
	"gitgud.io/softashell/rpgmaker-patch-translator/text"
	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
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

// Set the nasty global variables
func Init() {
	if stl == nil {
		stl = statictl.New()
	}
}

func ParseBlock(block PatchBlock) PatchBlock {
	if !text.ShouldTranslate(block.Original) {
		return block
	}

	block = ParseBlockLocalTL(block)
	block = ParseBlockRemoteTL(block)

	return block
}

func ParseBlockLocalTL(block PatchBlock) PatchBlock {
	var untranslated []string

	for i, t := range block.Translations {
		if t.Translated {
			continue // Block is already translated
		}

		// Attempt to get static translation for block
		blocks, untranslatedContexts, err := TranslateBlockStatic(t, block.Original)
		if err != nil || len(blocks) < 1 {
			continue
		}

		untranslated = append(untranslated, untranslatedContexts...)

		// Replace current
		if len(blocks) == 1 {
			block.Translations[i] = blocks[0]
		} else {
			// Delete current block
			block.Translations = append(block.Translations[:i], block.Translations[i+1:]...)

			// Append new blocks at end
			block.Translations = append(block.Translations, blocks...)
		}
	}

	// Leftovers
	if len(untranslated) > 0 {
		block.Translations = append(block.Translations, TranslationBlock{
			Text:       "",
			Contexts:   untranslated,
			Translated: false,
		})

		log.Infof("Mixed block in static\n %s", spew.Sdump(block))
	}

	return block
}

func ParseBlockRemoteTL(block PatchBlock) PatchBlock {
	var err error
	var items []lex.Item
	var untranslated []string
	var translated, parsed bool

	for i, t := range block.Translations {
		if t.Translated {
			continue // Block is already translated
		}

		// Fallback to lexing and translating chunks with comfy-translator
		good, bad := getTranslatableContexts(t, block.Original)
		untranslated = append(untranslated, bad...)

		if len(good) < 1 {
			continue
		}

		if !parsed {
			items, err = lex.ParseText(block.Original)
			if err != nil {
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

		log.Infof("Mixed block in comfy\n %s", spew.Sdump(block))
	}

	return block
}

func TranslateBlockStatic(b TranslationBlock, originalText string) ([]TranslationBlock, []string, error) {
	tlTypes := GetContextTypes(b.Contexts)
	blocks := []TranslationBlock{}
	untranslated := []string{}

	for t, c := range tlTypes {
		text, err := stl.GetTranslation(originalText, t)
		if err != nil {
			untranslated = append(untranslated, c...)
			continue
		}

		block := TranslationBlock{
			Text:       text,
			Contexts:   c,
			Touched:    true,
			Translated: true,
		}

		blocks = append(blocks, block)
	}

	return blocks, nil, nil
}

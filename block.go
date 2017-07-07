package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
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
	var items []item
	var err error

	parsed := false

	if !shouldTranslateText(block.original) {
		return block
	}

	for i, t := range block.translations {
		if shouldTranslateContext(t, block.original) {
			if !parsed {
				items, err = parseText(block.original)

				if err != nil {
					logBlockError(err, block, items)

					return block
				}

				parsed = true
			}

			t.text = translateItems(items)
			t.translated = true
			t.touched = true

			block.translations[i] = t

			log.Debugf("'%s' => '%s'\n", block.original, t.text)
		}
	}

	return block
}

func shouldTranslateContext(block translationBlock, text string) bool {
	if block.translated {
		log.Debug("Skipping translated block")
		return false
	}

	for _, c := range block.contexts {
		if engine == engineRPGMVX {
			//log.Debugf("%q", c)
			if strings.HasSuffix(c, "_se/name/") ||
				strings.HasSuffix(c, "/bgm/name/") ||
				strings.HasSuffix(c, "_me/name/") ||
				strings.Contains(c, "/InlineScript/") {
				return false
			}

			if strings.HasPrefix(c, ": Scripts/") {
				if strings.Contains(c, "Vocab/") {
					break
				} else {
					return false
				}
			}
		} else if engine == engineWolf {
			if strings.HasSuffix(c, "/Database") {
				return false
			} else if strings.HasPrefix(c, " DB:DataBase") {
				if strings.Contains(c, "アクター/") || //Actor
					strings.Contains(c, "NPC/") ||
					strings.Contains(c, "ステート/") || // State
					strings.Contains(c, "技能/") || // Skill
					strings.Contains(c, "敵/") || // Enemy
					strings.Contains(c, "武器/") || // Weapon
					strings.Contains(c, "称号/") || // Title
					strings.Contains(c, "衣装/") || // Clothing
					strings.Contains(c, "防具/") || // Armor
					strings.Contains(c, "道具/") || // Tools
					strings.Contains(c, "メニュー設計/") || // Menu
					strings.Contains(c, "コンフィグ/") || // Config
					strings.Contains(c, "クエスト/") || // Quest
					strings.Contains(c, "マップ選択画面") || // Map selection
					strings.Contains(c, "回想モード/") { // Recollection
					break
				}

				return false
			} else if strings.HasPrefix(c, " COMMONEVENT:") {
				if (strings.HasSuffix(c, "/SetString") && strings.Contains(text, "/")) || strings.HasSuffix(c, "/StringCondition") {
					return false
				}
			}
		}
	}

	return true
}

func translateItems(items []item) string {
	var out string

	for _, item := range items {
		log.Debugf("%14s: %q", item.typ, item.val)

		// Translate
		if item.typ == itemText {
			translation, err := translateString(item.val)
			check(err)

			if strings.HasPrefix(item.val, " ") && !strings.HasPrefix(translation, " ") {
				out += " "
			}

			out += translation

			if strings.HasSuffix(item.val, " ") && !strings.HasSuffix(translation, " ") {
				out += " "
			}
		} else if item.typ == itemEOF {
			continue
		} else if item.typ != itemError {
			// Add space before '(' since translation might make it get parsed as function
			if item.typ == itemRawString && item.val == "(" &&
				(!strings.HasSuffix(out, " ") && !strings.HasSuffix(out, "　")) {
				out += " "
			}

			// Add raw
			out += item.val
		}
	}

	return out
}

func parseText(text string) ([]item, error) {
	l := lex(text)

	var items []item
	var err error

	for {
		item := l.nextItem()
		items = append(items, item)

		if item.typ == itemError {
			err = fmt.Errorf("Failed to parse text")
		}

		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}

	return items, err
}

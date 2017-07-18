package main

import "strings"

func getTranslatableContexts(block translationBlock, text string) ([]string, []string) {
	var good, bad []string

	for _, c := range block.contexts {
		if shouldTranslateContext(c, text) {
			good = append(good, c)
		} else {
			bad = append(bad, c)
		}
	}

	return good, bad
}

func shouldTranslateContext(c, text string) bool {
	// TODO: Add switch to disable name translation to avoid breaking some games

	if engine == engineRPGMVX {
		//log.Debugf("%q", c)
		if strings.HasSuffix(c, "_se/name/") ||
			strings.HasSuffix(c, "/bgm/name/") ||
			strings.HasSuffix(c, "_me/name/") ||
			strings.Contains(c, "/InlineScript/") {
			return false
		}

		if strings.HasPrefix(c, ": Scripts/") {
			if strings.Contains(c, "Vocab/") ||
				(strings.HasPrefix(c, ": Scripts/Window_") &&
					(strings.Contains(c, "Info/") || strings.Contains(c, "Status/"))) {
				return true
			}

			return false
		}
	} else if engine == engineWolf {
		if strings.HasSuffix(c, "/Database") {
			return false
		} else if strings.HasPrefix(c, " DB:DataBase") {
			if strings.Contains(c, "アクター/") || //Actor
				strings.Contains(c, "キャラ名") || strings.Contains(c, "キャラクター名") || strings.HasSuffix(c, "名") || strings.HasSuffix(c, "名前") || // Character name
				strings.Contains(c, "タイトル") || // Title
				strings.Contains(c, "NPC/") ||
				strings.Contains(c, "ステート/") || strings.Contains(c, "状態名") || // State
				strings.Contains(c, "技能/") || // Skill
				strings.Contains(c, "敵/") || // Enemy
				strings.Contains(c, "武器/") || // Weapon
				strings.Contains(c, "称号/") || // Title
				strings.Contains(c, "衣装/") || // Clothing
				strings.Contains(c, "防具/") || // Armor
				strings.Contains(c, "道具/") || // Tools
				strings.Contains(c, "メニュー設計/") || // Menu
				strings.Contains(c, "戦闘コマンド/") || // Battle
				strings.Contains(c, "コンフィグ/") || strings.Contains(c, "用語設定/") || // Config
				strings.Contains(c, "クエスト/") || // Quest
				strings.Contains(c, "依頼主") || // Client name
				strings.Contains(c, "マップ選択画面") || // Map selection
				strings.Contains(c, "回想モード/") { // Recollection
				return true
			}

			return false
		} else if strings.HasPrefix(c, " COMMONEVENT:") {
			if (strings.HasSuffix(c, "/SetString") && strings.Contains(text, "/")) || strings.HasSuffix(c, "/StringCondition") {
				return false
			}
		} else if strings.HasPrefix(c, " GAMEDAT:") && !strings.HasSuffix(c, "Title") {
			return false
		}
	}

	return true
}

func shouldBreakLines(contexts []string) bool {
	for _, c := range contexts {
		if engine == engineRPGMVX {
			if strings.Contains(c, "GameINI/Title") || strings.Contains(c, "System/game_title/") {
				return false
			}
		} else if engine == engineWolf {
			if strings.HasPrefix(c, " GAMEDAT:") {
				return false
			}
		}
	}

	return true
}

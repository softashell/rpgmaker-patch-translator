package block

import (
	"strings"

	"gitgud.io/softashell/rpgmaker-patch-translator/statictl"
)

func GetTranslationType(contexts []string) statictl.TranslationType {
	// TODO: Handle multiple contexts instead of picking first one that matches

	for _, c := range contexts {

		if IsActor(c) || IsArmor(c) || IsClass(c) || IsEnemy(c) || IsItem(c) || IsSkill(c) || IsTroop(c) || IsWeapon(c) {
			if IsName(c) {
				return statictl.TransName
			}

			if IsDescription(c) {
				return statictl.TransDescription
			}
		}

		if IsDialogue(c) {
			return statictl.TransDialogue
		}

		if IsChoice(c) {
			return statictl.TransChoice
		}

		if IsMessage(c) {
			return statictl.TransMessage
		}

		if IsInlineScript(c) {
			return statictl.TransInlineScript
		}

		if IsScript(c) {
			if IsVocab(c) {
				return statictl.TransVocab
			}

			return statictl.TransScript
		}

		if IsSystem(c) {
			return statictl.TransSystem
		}
	}

	return statictl.TransGeneric
}

func IsMessage(c string) bool {
	if strings.Contains(c, "/message") {
		return true
	}

	return false
}

func IsDialogue(c string) bool {
	if strings.HasSuffix(c, "/Dialogue") {
		return true
	}

	return false
}

func IsChoice(c string) bool {
	if strings.Contains(c, "/Choice/") {
		return true
	}

	return false
}

func IsVocab(c string) bool {
	if strings.Contains(c, "/Vocab") {
		return true
	}

	return false
}

func IsName(c string) bool {
	if strings.HasSuffix(c, "/name/") {
		return true
	}

	return false
}

func IsDescription(c string) bool {
	if strings.HasSuffix(c, "/description//") {
		return true
	}

	return false
}

func IsActor(c string) bool {
	if strings.HasPrefix(c, ": Actors/") {
		return true
	}

	return false
}

func IsEnemy(c string) bool {
	if strings.HasPrefix(c, ": Enemies/") {
		return true
	}

	return false
}

func IsTroop(c string) bool {
	if strings.HasPrefix(c, ": Troops/") {
		return true
	}

	return false
}

func IsArmor(c string) bool {
	if strings.HasPrefix(c, ": Armors/") {
		return true
	}

	return false
}

func IsWeapon(c string) bool {
	if strings.HasPrefix(c, ": Weapons/") {
		return true
	}

	return false
}

func IsItem(c string) bool {
	if strings.HasPrefix(c, ": Items/") {
		return true
	}

	return false
}

func IsClass(c string) bool {
	if strings.HasPrefix(c, ": Classes/") {
		return true
	}

	return false
}

func IsState(c string) bool {
	if strings.HasPrefix(c, ": States/") {
		return true
	}

	return false
}

func IsSkill(c string) bool {
	if strings.HasPrefix(c, ": Skills/") {
		return true
	}

	return false
}

func IsInlineScript(c string) bool {
	if strings.Contains(c, "/InlineScript/") {
		return true
	}

	return false
}

func IsSystem(c string) bool {
	if strings.HasPrefix(c, ": System/") {
		return true
	}

	return false
}

func IsScript(c string) bool {
	if strings.HasPrefix(c, ": Scripts/") {
		return true
	}

	return false
}

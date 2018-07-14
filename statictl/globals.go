package statictl

type TranslationType int

const (
	TransGeneric     TranslationType = iota // Default group
	TransName                               // Every name
	TransDescription                        // Every description
	TransDialogue
	TransChoice
	TransVocab
	TransMessage      // Skill and Status messages
	TransInlineScript // Inline scripts
	TransScript       // Scripts
	TransSystem       // System translations like main menu and title
)

var databaseFiles = map[TranslationType]string{
	TransGeneric:      "Generic.hjson",
	TransName:         "Name.hjson",
	TransDescription:  "Description.hjson",
	TransDialogue:     "Dialogue.hjson",
	TransChoice:       "Choice.hjson",
	TransVocab:        "Vocab.hjson",
	TransMessage:      "Message.hjson",
	TransInlineScript: "InlineScript.hjson",
	TransScript:       "Script.hjson",
	TransSystem:       "System.hjson",
}

type DatabaseType int

const (
	DbStatic DatabaseType = iota
	DbDynamic
)

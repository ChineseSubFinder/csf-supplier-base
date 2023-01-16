package settings

type FailWordsConfig struct {
	Enable     bool
	Words      []string
	WordsRegex []string
}

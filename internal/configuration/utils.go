package configuration

func IsSource(name string) bool {
	return name == "source"
}

// IsLanguageSupportedWithOptions returns true if the specified language supports options.
func IsLanguageSupportedWithOptions(language string) bool {
	return supportedLanguagesWithOptions[language]
}

// GetSupportedLanguagesWithOptions returns a slice of languages that support options.
func GetSupportedLanguagesWithOptions() []string {
	languages := make([]string, 0, len(supportedLanguagesWithOptions))
	for lang := range supportedLanguagesWithOptions {
		languages = append(languages, lang)
	}
	return languages
}

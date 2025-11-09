// package language defines the supported programming languages and provides
// metadata associated with each one, such as file extensions and comment styles.
package language

// Language is a string type representing a supported programming language.
type Language string

// These constants enumerate the programming languages supported by the application.
const (
	Cpp     Language = "cpp"
	Go      Language = "go"
	Java    Language = "java"
	Kotlin  Language = "kotlin"
	Lua     Language = "lua"
	Php     Language = "php"
	Swift   Language = "swift"
	Dart    Language = "dart"
	Csharp  Language = "csharp"
	Python  Language = "python"
	Rust    Language = "rust"
	Ts      Language = "ts"
	Nim     Language = "nim"
	Unknown Language = "unknown"
)

// LanguageMetadata holds information specific to a programming language,
// such as its typical file extension and syntax for line comments.
type LanguageMetadata struct {
	Language     Language
	CommentStyle string
	Extension    string
}

// languageMetadataMap provides a central mapping from a Language to its metadata.
var languageMetadataMap = map[Language]LanguageMetadata{
	Cpp:     {Language: Cpp, CommentStyle: "//", Extension: ".h"},
	Go:      {Language: Go, CommentStyle: "//", Extension: ".go"},
	Java:    {Language: Java, CommentStyle: "//", Extension: ".java"},
	Kotlin:  {Language: Kotlin, CommentStyle: "//", Extension: ".kt"},
	Lua:     {Language: Lua, CommentStyle: "--", Extension: ".lua"},
	Php:     {Language: Php, CommentStyle: "//", Extension: ".php"},
	Swift:   {Language: Swift, CommentStyle: "//", Extension: ".swift"},
	Dart:    {Language: Dart, CommentStyle: "//", Extension: ".dart"},
	Csharp:  {Language: Csharp, CommentStyle: "//", Extension: ".cs"},
	Python:  {Language: Python, CommentStyle: "#", Extension: ".py"},
	Rust:    {Language: Rust, CommentStyle: "//", Extension: ".rs"},
	Ts:      {Language: Ts, CommentStyle: "//", Extension: ".ts"},
	Nim:     {Language: Nim, CommentStyle: "#", Extension: ".nim"},
	Unknown: {Language: Unknown, CommentStyle: "", Extension: ""},
}

// GetMetadata retrieves the metadata for a given language. It returns an
// UnsupportedLanguageError if the language is not found in the map.
func GetMetadata(lang Language) (LanguageMetadata, error) {
	metadata, exists := languageMetadataMap[lang]
	if !exists {
		return metadata, &UnsupportedLanguageError{}
	}
	return metadata, nil
}

// GetSupportedLanguages returns a slice containing all supported languages,
// excluding the 'Unknown' type.
func GetSupportedLanguages() []Language {
	languages := make([]Language, 0, len(languageMetadataMap))
	for lang := range languageMetadataMap {
		if lang != Unknown {
			languages = append(languages, lang)
		}
	}
	return languages
}

// IsSupportedLanguage checks if the given language string corresponds to a
// supported language.
func IsSupportedLanguage(language string) bool {
	_, ok := languageMetadataMap[Language(language)]
	return ok
}

// UnsupportedLanguageError is returned when an operation is attempted on a
// language that is not supported by the application.
type UnsupportedLanguageError struct{}

// Error implements the error interface for UnsupportedLanguageError.
func (e *UnsupportedLanguageError) Error() string {
	return "unsupported language: this language is not yet supported"
}

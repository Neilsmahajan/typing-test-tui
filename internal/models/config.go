package models

import "strings"

type Mode string

const (
	QuoteMode Mode = "quote"
	WordsMode Mode = "words"
	TimeMode  Mode = "time"
)

type Language string

const (
	Chinese    Language = "chinese_simplified"
	Assembly   Language = "code_assembly"
	C          Language = "code_c"
	Cpp        Language = "code_c++"
	CSharp     Language = "code_csharp"
	CSS        Language = "code_css"
	Go         Language = "code_go"
	Java       Language = "code_java"
	JavaScript Language = "code_javascript"
	Kotlin     Language = "code_kotlin"
	Lua        Language = "code_lua"
	PHP        Language = "code_php"
	Python     Language = "code_python"
	R          Language = "code_r"
	Ruby       Language = "code_ruby"
	Rust       Language = "code_rust"
	Typescript Language = "code_typescript"
	English    Language = "english"
	French     Language = "french"
	Spanish    Language = "spanish"
)

type Duration int

type WordCount int

type Config struct {
	Mode      Mode
	Language  Language
	Duration  Duration
	WordCount WordCount
}

var supportedLanguages = []Language{
	Chinese,
	Assembly,
	C,
	Cpp,
	CSharp,
	CSS,
	Go,
	Java,
	JavaScript,
	Kotlin,
	Lua,
	PHP,
	Python,
	R,
	Ruby,
	Rust,
	Typescript,
	English,
	French,
	Spanish,
}

var languageAliases = map[string]Language{
	"chinese":            Chinese,
	"zh":                 Chinese,
	"zh_cn":              Chinese,
	"chinese_simplified": Chinese,
	"assembly":           Assembly,
	"asm":                Assembly,
	"code_assembly":      Assembly,
	"c":                  C,
	"code_c":             C,
	"cpp":                Cpp,
	"c++":                Cpp,
	"code_cpp":           Cpp,
	"csharp":             CSharp,
	"c#":                 CSharp,
	"code_csharp":        CSharp,
	"css":                CSS,
	"code_css":           CSS,
	"go":                 Go,
	"golang":             Go,
	"code_go":            Go,
	"java":               Java,
	"code_java":          Java,
	"javascript":         JavaScript,
	"js":                 JavaScript,
	"code_javascript":    JavaScript,
	"kotlin":             Kotlin,
	"code_kotlin":        Kotlin,
	"lua":                Lua,
	"code_lua":           Lua,
	"php":                PHP,
	"code_php":           PHP,
	"python":             Python,
	"py":                 Python,
	"code_python":        Python,
	"r":                  R,
	"code_r":             R,
	"ruby":               Ruby,
	"code_ruby":          Ruby,
	"rust":               Rust,
	"code_rust":          Rust,
	"typescript":         Typescript,
	"ts":                 Typescript,
	"code_typescript":    Typescript,
	"english":            English,
	"en":                 English,
	"eng":                English,
	"spanish":            Spanish,
	"es":                 Spanish,
	"french":             French,
	"fr":                 French,
}

func NormalizeLanguage(input string) (Language, bool) {
	key := strings.TrimSpace(strings.ToLower(input))
	key = strings.ReplaceAll(key, "-", "_")
	if lang, ok := languageAliases[key]; ok {
		return lang, true
	}

	for _, lang := range supportedLanguages {
		if key == string(lang) {
			return lang, true
		}
	}

	return Language(key), false
}

func SupportedLanguages() []Language {
	result := make([]Language, len(supportedLanguages))
	copy(result, supportedLanguages)
	return result
}

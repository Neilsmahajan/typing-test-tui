package models

type Mode string

const (
	QuoteMode Mode = "quote"
	WordsMode Mode = "words"
	TimeMode  Mode = "time"
)

type Language string

const (
	Chinese    Language = "chinese"
	Assembly   Language = "assembly"
	C          Language = "c"
	Cpp        Language = "cpp"
	CSharp     Language = "csharp"
	CSS        Language = "css"
	Go         Language = "go"
	Java       Language = "java"
	JavaScript Language = "javascript"
	Kotlin     Language = "kotlin"
	Lua        Language = "lua"
	PHP        Language = "php"
	Python     Language = "python"
	R          Language = "r"
	Ruby       Language = "ruby"
	Rust       Language = "rust"
	Typescript Language = "typescript"
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
package typing

import (
	"math/rand"
	"strings"

	"github.com/neilsmahajan/typing-test-tui/internal/models"
)

// PunctuateWord applies punctuation to a word based on language and context.
// This is a simplified implementation based on Monkeytype's punctuation logic.
func PunctuateWord(
	rng *rand.Rand,
	language models.Language,
	previousWord string,
	currentWord string,
	index int,
	maxIndex int,
) string {
	if rng == nil || currentWord == "" {
		return currentWord
	}

	word := currentWord
	langPrefix := getLanguagePrefix(language)
	lastChar := getLastChar(previousWord)

	// First word or word after sentence-ending punctuation: capitalize
	if langPrefix != "code" && (index == 0 || shouldCapitalize(lastChar)) {
		word = capitalizeFirstLetter(word)
	}

	// Sentence-ending punctuation (10% chance, or always on last word)
	if (rng.Float64() < 0.1 && lastChar != "." && lastChar != "," && index != maxIndex-2) || index == maxIndex-1 {
		word = addSentenceEnding(rng, word, langPrefix)
	} else if rng.Float64() < 0.2 && lastChar != "," {
		// Comma (20% chance)
		word = addComma(word, langPrefix)
	} else if rng.Float64() < 0.01 && lastChar != "," && lastChar != "." {
		// Quotes (1% chance)
		word = `"` + word + `"`
	} else if rng.Float64() < 0.011 && lastChar != "," && lastChar != "." {
		// Single quotes (1.1% chance)
		word = "'" + word + "'"
	} else if rng.Float64() < 0.012 && lastChar != "," && lastChar != "." {
		// Parentheses (1.2% chance)
		word = addParentheses(rng, word, langPrefix)
	} else if rng.Float64() < 0.013 && lastChar != "," && lastChar != "." && lastChar != ";" && lastChar != ":" {
		// Colon (1.3% chance)
		word = addColon(word, langPrefix)
	} else if rng.Float64() < 0.014 && lastChar != "," && lastChar != "." && previousWord != "-" {
		// Dash (1.4% chance)
		word = "-"
	} else if rng.Float64() < 0.015 && lastChar != "," && lastChar != "." && lastChar != ";" {
		// Semicolon (1.5% chance)
		word = addSemicolon(word, langPrefix)
	} else if rng.Float64() < 0.25 && langPrefix == "code" {
		// Code-specific punctuation (25% chance)
		word = addCodePunctuation(rng, word, language)
	}

	return word
}

// getLanguagePrefix extracts the language prefix (e.g., "code", "english", "spanish")
func getLanguagePrefix(language models.Language) string {
	langStr := string(language)
	if strings.HasPrefix(langStr, "code_") {
		return "code"
	}
	return strings.Split(langStr, "_")[0]
}

// getLastChar returns the last character of a string
func getLastChar(s string) string {
	if len(s) == 0 {
		return ""
	}
	runes := []rune(s)
	return string(runes[len(runes)-1])
}

// shouldCapitalize returns true if the character indicates a new sentence
func shouldCapitalize(lastChar string) bool {
	return lastChar == "." || lastChar == "?" || lastChar == "!" || lastChar == "。" || lastChar == "।"
}

// capitalizeFirstLetter capitalizes the first letter of a word
func capitalizeFirstLetter(word string) string {
	if len(word) == 0 {
		return word
	}
	runes := []rune(word)
	if runes[0] >= 'a' && runes[0] <= 'z' {
		runes[0] = runes[0] - 'a' + 'A'
	}
	return string(runes)
}

// addSentenceEnding adds a sentence-ending punctuation mark
func addSentenceEnding(rng *rand.Rand, word string, langPrefix string) string {
	rand := rng.Float64()

	if langPrefix == "chinese" {
		if rand <= 0.8 {
			return word + "。"
		} else if rand <= 0.9 {
			return word + "？"
		}
		return word + "！"
	}

	if rand <= 0.8 {
		return word + "."
	} else if rand <= 0.9 {
		return word + "?"
	}
	return word + "!"
}

// addComma adds a comma to the word
func addComma(word string, langPrefix string) string {
	switch langPrefix {
	case "chinese":
		return word + "，"
	default:
		return word + ","
	}
}

// addParentheses adds parentheses around the word
func addParentheses(rng *rand.Rand, word string, langPrefix string) string {
	if langPrefix == "code" {
		brackets := []string{"()", "{}", "[]", "<>"}
		bracket := brackets[rng.Intn(len(brackets))]
		return string(bracket[0]) + word + string(bracket[1])
	}
	if langPrefix == "chinese" {
		return "（" + word + "）"
	}
	return "(" + word + ")"
}

// addColon adds a colon to the word
func addColon(word string, langPrefix string) string {
	if langPrefix == "chinese" {
		return word + "："
	}
	return word + ":"
}

// addSemicolon adds a semicolon to the word
func addSemicolon(word string, langPrefix string) string {
	if langPrefix == "chinese" {
		return word + "；"
	}
	return word + ";"
}

// addCodePunctuation adds code-specific punctuation
func addCodePunctuation(rng *rand.Rand, word string, language models.Language) string {
	langStr := string(language)

	// C-family languages get extended operators
	if strings.HasPrefix(langStr, "code_c") && !strings.HasPrefix(langStr, "code_css") {
		specials := []string{
			"{", "}", "[", "]", "(", ")", ";", "=", "+", "%", "/",
			"/*", "*/", "//", "!=", "==", "<=", ">=", "||", "&&",
			"<<", ">>", "%=", "&=", "*=", "++", "+=", "--", "-=",
			"/=", "^=", "|=",
		}
		return specials[rng.Intn(len(specials))]
	}

	// JavaScript gets backticks
	if langStr == "code_javascript" || langStr == "code_typescript" {
		specials := []string{"{", "}", "[", "]", "(", ")", ";", "=", "+", "%", "/", "`"}
		return specials[rng.Intn(len(specials))]
	}

	// Default code punctuation
	specials := []string{"{", "}", "[", "]", "(", ")", ";", "=", "+", "%", "/"}
	return specials[rng.Intn(len(specials))]
}

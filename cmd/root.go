package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/neilsmahajan/typing-test-tui/internal/app"
	"github.com/neilsmahajan/typing-test-tui/internal/models"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "typing-test-tui",
	Short: "A typing test application in the terminal",
	Long: `The typing-test-tui is a terminal-based typing test application built with Go and Bubble Tea.
It provides an interactive interface to practice and improve your typing skills.`,
	Run: runTypingTest,
}

const (
	defaultDuration  = 60
	defaultWordCount = 50
)

var (
	allowedDurations    = []int{15, 30, 60, 120}
	allowedDurationSet  = map[int]struct{}{15: {}, 30: {}, 60: {}, 120: {}}
	allowedWordCounts   = []int{10, 25, 50, 100}
	allowedWordCountSet = map[int]struct{}{10: {}, 25: {}, 50: {}, 100: {}}
)

func runTypingTest(cmd *cobra.Command, args []string) {
	mode, err := cmd.Flags().GetString("mode")
	if err != nil {
		fmt.Println("Error reading mode flag:", err)
		return
	}

	language, err := cmd.Flags().GetString("language")
	if err != nil {
		fmt.Println("Error reading language flag:", err)
		return
	}

	duration, err := cmd.Flags().GetInt("duration")
	if err != nil {
		fmt.Println("Error reading duration flag:", err)
		return
	}

	wordCount, err := cmd.Flags().GetInt("word-count")
	if err != nil {
		fmt.Println("Error reading word count flag:", err)
		return
	}

	modeValue := models.Mode(mode)

	if err := validateFlags(modeValue, duration, wordCount); err != nil {
		fmt.Println("Error:", err)
		return
	}

	normalizedLanguage, err := normalizeLanguage(language)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	cfg := models.Config{
		Mode:      modeValue,
		Language:  normalizedLanguage,
		Duration:  models.Duration(duration),
		WordCount: models.WordCount(wordCount),
	}

	if err := app.Run(cfg); err != nil {
		fmt.Println("Error running app:", err)
	}
}

func validateFlags(mode models.Mode, duration int, wordCount int) error {
	switch mode {
	case models.QuoteMode:
		if duration != defaultDuration {
			return fmt.Errorf("duration flag is only available for time mode")
		}
		if wordCount != defaultWordCount {
			return fmt.Errorf("word-count flag is only available for words mode")
		}
	case models.WordsMode:
		if duration != defaultDuration {
			return fmt.Errorf("duration flag is only available for time mode")
		}
		if _, ok := allowedWordCountSet[wordCount]; !ok {
			return fmt.Errorf("word count must be one of %s", joinInts(allowedWordCounts))
		}
	case models.TimeMode:
		if _, ok := allowedDurationSet[duration]; !ok {
			return fmt.Errorf("duration must be one of %s", joinInts(allowedDurations))
		}
		if wordCount != defaultWordCount {
			return fmt.Errorf("word-count flag is only available for words mode")
		}
	default:
		return fmt.Errorf("unsupported mode %q. Supported modes: 'quote', 'words', 'time'", mode)
	}

	return nil
}

func normalizeLanguage(language string) (models.Language, error) {
	if lang, ok := models.NormalizeLanguage(language); ok {
		return lang, nil
	}

	supported := models.SupportedLanguages()
	names := make([]string, len(supported))
	for i, lang := range supported {
		names[i] = string(lang)
	}

	return models.Language(""), fmt.Errorf("unsupported language %q. Supported languages: %s", language, strings.Join(names, ", "))
}

func joinInts(values []int) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, ", ")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.typing-test-tui.yaml)")
	rootCmd.Flags().StringP("mode", "m", "quote", "Mode of the typing test ('quote', 'words', 'time')")
	rootCmd.Flags().StringP("language", "l", "english", "Language for the typing test (e.g., 'english' for English, 'spanish' for Spanish, 'go' for Go code)")
	rootCmd.Flags().IntP("duration", "d", 60, "Duration of the typing test in seconds (only for 'time' mode; options: 15, 30, 60, 120)")
	rootCmd.Flags().IntP("word-count", "w", 50, "Number of words for the typing test (only for 'words' mode; options: 10, 25, 50, 100)")
}

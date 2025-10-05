package cmd

import (
	"fmt"
	"os"

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

	var cfg = models.Config{
		Mode:      models.Mode(mode),
		Language:  models.Language(language),
		Duration:  models.Duration(duration),
		WordCount: models.WordCount(wordCount),
	}

	if err := app.Run(cfg); err != nil {
		fmt.Println("Error running app:", err)
		return
	}
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
	rootCmd.Flags().IntP("duration", "d", 60, "Duration of the typing test in seconds (only for 'time' mode)")
	rootCmd.Flags().IntP("word-count", "w", 50, "Number of words for the typing test (only for 'words' mode)")
}

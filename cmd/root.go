package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/neilsmahajan/typing-test-tui/internal/ui/quote_input"
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

	var p *tea.Program

	switch mode {
	case "quote":
		p = tea.NewProgram(quote_input.InitialModel("The quick brown fox jumps over the lazy dog."))
	default:
		fmt.Println("Unsupported mode:", mode)
		return
	}

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.typing-test-tui.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("mode", "m", "quote", "Mode of the typing test ('quote', 'words', 'time')")
}

package cmd

import "github.com/spf13/cobra"

var modesCmd = &cobra.Command{
	Use:     "modes",
	Aliases: []string{"list-modes"},
	Short:   "List supported modes",
	Long:    `List all the modes supported by the typing test application.`,
	Example: "typing-test-tui modes",
	Args:    cobra.NoArgs,
	Run:     listModes,
}

func listModes(cmd *cobra.Command, args []string) {
	cmd.Println("Supported Modes:")
	cmd.Println(" - quote : Type predefined quotes.")
	cmd.Println(" - words : Type a set number of random words.")
	cmd.Println(" - time  : Type as many words as you can in a set time limit.")
	cmd.Println("\nYou can specify a mode using the --mode or -m flag when starting a typing test.")
}

func init() {
	rootCmd.AddCommand(modesCmd)
}

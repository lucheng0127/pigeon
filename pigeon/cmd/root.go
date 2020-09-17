package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pigeon",
	Short: "Pigeon a toolkit to manage hosts and scripts",
	Long:  "Pigeon a toolkit to manage hosts and scripts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run pigeon -h for more information")
	},
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Registry commands
	rootCmd.AddCommand(scriptListCmd)
	rootCmd.AddCommand(scriptAddCmd)
	rootCmd.AddCommand(versionCmd)
}

// Execute cmd
func Execute() {
	// Run root command
	err := rootCmd.Execute()
	checkError(err)
}

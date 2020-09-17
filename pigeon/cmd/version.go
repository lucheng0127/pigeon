package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of Pigeon",
	Long:  `Version of Pigeon`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Pigeon 0.0.1")
	},
}

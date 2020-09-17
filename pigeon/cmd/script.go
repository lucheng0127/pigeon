package cmd

import (
	"fmt"
	"pigeon/modules/sockets"

	"github.com/spf13/cobra"
)

// List scritp command flags
var name string
var scriptFile string
var pigeonSocket = sockets.UnixSocket{SocketFile: "/var/run/pigeond.socket"}

var scriptListCmd = &cobra.Command{
	Use:   "script-list",
	Short: "List scripts",
	Long:  "List scripts",
	Run: func(cmd *cobra.Command, args []string) {
		// Send list script command to unix socket and get result
		rst := sockets.Send(&pigeonSocket, "F LIST_SCRIPTS END")
		rstData, _ := checkJSONRst(rst)
		fmt.Println(rstData["Result"])
	},
}

var scriptAddCmd = &cobra.Command{
	Use:   "script-add",
	Short: "Add script",
	Long:  "Add script into script inventory",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Add script %s with script file %s", name, scriptFile)
	},
}

func init() {
	// Add flags, some of it required
	scriptAddCmd.Flags().StringVarP(&name, "name", "n", "", "Script name (required)")
	scriptAddCmd.MarkFlagRequired("name")
	scriptAddCmd.Flags().StringVarP(&scriptFile, "file", "f", "", "Script file (required)")
	scriptAddCmd.MarkFlagRequired("file")
}

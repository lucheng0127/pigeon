package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"pigeon/modules/sockets"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// List scritp command flags
var name string
var scriptFile string
var scriptFileType string
var passwd string
var pigeonSocket = sockets.UnixSocket{SocketFile: "/var/run/pigeon/pigeond.socket"}

func loadCSV(rawData string) [][]string {
	lines, err := csv.NewReader(strings.NewReader(rawData)).ReadAll()
	checkError(err)
	return lines
}

var scriptListCmd = &cobra.Command{
	Use:   "script-list",
	Short: "List scripts",
	Long:  "List scripts",
	Run: func(cmd *cobra.Command, args []string) {
		// Send list script command to unix socket and get result
		rst := sockets.Send(&pigeonSocket, "F LIST_SCRIPTS END")
		rstData, _ := checkJSONRst(rst)
		rawScriptsData := fmt.Sprint(rstData["Result"])
		lines := loadCSV(rawScriptsData)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Script Name", "Create Time", "MD5SUM"})
		for _, line := range lines {
			cTimestamp, err := strconv.ParseInt(line[1], 10, 64)
			checkError(err)
			cTime := time.Unix(cTimestamp, 0).Format("2006-01-02 15:04:05")
			table.Append([]string{line[0], cTime, line[2]})
		}
		table.Render()
	},
}

var scriptAddCmd = &cobra.Command{
	Use:   "script-add",
	Short: "Add script",
	Long:  "Add script into script inventory",
	Run: func(cmd *cobra.Command, args []string) {
		// Replace space to + in name
		name = strings.ReplaceAll(name, " ", "+")

		cmdStr := fmt.Sprintf("F ADD_SCRIPT %s %s", name, scriptFile)
		if scriptFileType != "" {
			cmdStr += fmt.Sprintf(" %s", scriptFileType)
		} else {
			cmdStr += " +"
		}
		if passwd != "" {
			cmdStr += fmt.Sprintf(" %s", passwd)
		} else {
			cmdStr += " +"
		}
		cmdStr += " END"
		rst := sockets.Send(&pigeonSocket, cmdStr)
		rstData, _ := checkJSONRst(rst)
		ec := fmt.Sprint(rstData["ExitCode"])
		if ec == "1" {
			fmt.Print("Failed: ")
		}
		fmt.Println(rstData["Result"])
	},
}

func init() {
	// Add flags, some of it required
	scriptAddCmd.Flags().StringVarP(&name, "name", "n", "", "Script name (required)")
	scriptAddCmd.MarkFlagRequired("name")
	scriptAddCmd.Flags().StringVarP(&scriptFile, "file", "f", "", "Script file (required)")
	scriptAddCmd.MarkFlagRequired("file")
	scriptAddCmd.Flags().StringVarP(&scriptFileType, "type", "t", "", "Script compress type tar/zip")
	scriptAddCmd.Flags().StringVarP(&passwd, "passwd", "p", "", "Script compress password")
}

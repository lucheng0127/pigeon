package tasks

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"pigeon/pigeond/log"
	"strconv"
)

const scriptsInventoryFile = "/var/run/pigeon/scirpt_inventory.csv"
const scriptDir = "/var/run/pigeon/scripts" // Script store root
var scriptsInventoryData []script

// script define struct of scripts
type script struct {
	ScriptName   string // Script name
	CreateTime   int64  // Create time
	Md5sum       string // MD5SUM of scripts file
	CompressType string // Compress type tar or zip
	Encrypted    bool   // Is encrypted
	Password     string // Encrypt password
}

type scriptsInventory struct {
	scripts []script // Scripts list
}

func (s *script) addToCSV() {
	file, err := os.OpenFile(scriptsInventoryFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	checkErr(err)
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()
	for _, s := range scriptsInventoryData {
		err := w.Write([]string{
			s.ScriptName,
			strconv.Itoa(int(s.CreateTime)),
			s.Md5sum,
			strconv.FormatBool(s.Encrypted),
			s.Password,
		})
		checkErr(err)
	}
}

func loadCSV() {
	file, err := os.OpenFile(scriptsInventoryFile, os.O_RDONLY, 0666)
	checkErr(err)
	defer file.Close()
	lines, err := csv.NewReader(file).ReadAll()
	checkErr(err)
	s := script{}
	for _, line := range lines {
		s.ScriptName = line[0]
		s.CreateTime, err = strconv.ParseInt(line[1], 10, 64)
		checkErr(err)
		s.Md5sum = line[2]
		s.CompressType = line[3]
		s.Encrypted, err = strconv.ParseBool(line[4])
		checkErr(err)
		s.Password = line[5]
		scriptsInventoryData = append(scriptsInventoryData, s)
		s = script{}
	}

}

func init() {
	file, err := os.OpenFile(scriptsInventoryFile, os.O_RDONLY|os.O_CREATE, 0666)
	checkErr(err)
	file.Close()
	// Load scripts inventory data into scriptsInventoryData
	loadCSV()
	log.Log.Debug("Load scripts inventory finished")
}

func checkErr(err error) {
	if err != nil {
		log.Log.Error(err.Error())
	}
}

// ListScript to list script list
func ListScript() (int, string) {
	f, err := os.OpenFile(scriptsInventoryFile, os.O_RDONLY, 0666)
	f.Close()
	checkErr(err)
	scriptsData, err := ioutil.ReadFile(scriptsInventoryFile)
	checkErr(err)
	return 0, string(scriptsData)
}

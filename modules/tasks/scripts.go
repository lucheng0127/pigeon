package tasks

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"pigeon/pigeond/log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const scriptsInventoryFile = "/var/run/pigeon/script_inventory.csv"
const scriptDir = "/var/run/pigeon/scripts" // Script store root
var scriptsInventoryData []script
var muLock sync.Mutex = sync.Mutex{}
var scriptKeyMap = make(map[string]bool)

// script define struct of scripts
type script struct {
	ScriptName   string // Script name
	CreateTime   int64  // Create time
	Md5sum       string // MD5SUM of scripts file
	CompressType string // Compress type tar or zip
	Password     string // Encrypt password
	Encrypted    bool   // Is encrypted
}

type scriptsInventory struct {
	scripts []script // Scripts list
}

func existedKey(key string) bool {
	return scriptKeyMap["key"]
}

func (s script) addToCSV() (int, string) {

	// Open script inventory csv file with append module
	file, err := os.OpenFile(scriptsInventoryFile, os.O_APPEND|os.O_RDWR, 0666)
	checkErr(err)
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()

	// Add script into csv and script inventory
	muLock.Lock()
	err = w.Write([]string{
		s.ScriptName,
		strconv.Itoa(int(s.CreateTime)),
		s.Md5sum,
		s.CompressType,
		s.Password,
		strconv.FormatBool(s.Encrypted),
	})
	if err != nil {
		log.Log.Error(err.Error())
		return 1, err.Error()
	}

	scriptsInventoryData = append(scriptsInventoryData, s)
	scriptKeyMap[s.ScriptName] = true
	muLock.Unlock()
	return 0, "Add script succed"
}

func loadCSV() {
	file, err := os.OpenFile(scriptsInventoryFile, os.O_RDONLY, 0666)
	checkErr(err)
	defer file.Close()
	lines, err := csv.NewReader(file).ReadAll()
	checkErr(err)
	s := script{}
	for _, line := range lines {

		// Load script data from csv
		s.ScriptName = line[0]
		s.CreateTime, err = strconv.ParseInt(line[1], 10, 64)
		checkErr(err)
		s.Md5sum = line[2]
		s.CompressType = line[3]
		s.Password = line[4]
		s.Encrypted, err = strconv.ParseBool(line[5])
		checkErr(err)

		// Add lock when append to script inventory
		muLock.Lock()
		scriptsInventoryData = append(scriptsInventoryData, s)
		scriptKeyMap[line[0]] = true
		muLock.Unlock()

		// Clean script struct
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

// To fix move file between different drive
func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

func hashFile(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Log.Error(err.Error())
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Log.Error(err.Error())
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
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

// AddScript add a script into script inventory
func AddScript(name, file, fType, passwd string) (int, string) {

	if name == "" {
		return 1, "Script name is required"
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return 1, "Script file is required"
	}

	// Check script name
	if existedKey(name) {
		return 1, "Script Name not unique"
	}

	// Copy script file and calculate hash
	s := script{
		ScriptName: name,
		CreateTime: time.Now().Unix(),
	}
	if passwd != "" {
		s.Encrypted = true
	} else {
		s.Password = "*"
		s.Encrypted = false
	}
	switch strings.ToUpper(fType) {
	case "ZIP":
		s.CompressType = "ZIP"
	default:
		s.CompressType = "TAR"
	}
	if _, err := os.Stat(scriptDir); os.IsNotExist(err) {
		err := os.Mkdir(scriptDir, 0755)
		if err != nil {
			return 1, err.Error()
		}
	}
	scritpStoreDir := path.Join(scriptDir, strconv.Itoa(int(s.CreateTime)))
	if _, err := os.Stat(scritpStoreDir); os.IsNotExist(err) {
		err := os.Mkdir(scritpStoreDir, 0755)
		if err != nil {
			return 1, err.Error()
		}
	}
	scriptFileName := path.Join(scritpStoreDir, "script."+strings.ToLower(s.CompressType))
	err := moveFile(file, scriptFileName)
	if err != nil {
		os.RemoveAll(scritpStoreDir)
		return 1, err.Error()
	}
	// Caculate md5
	fileHash, err := hashFile(scriptFileName)
	if err != nil {
		return 1, err.Error()
	}
	s.Md5sum = fileHash

	return s.addToCSV()
}

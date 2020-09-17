package cmd

import (
	"encoding/json"
)

func checkJSONRst(rst string) (map[string]interface{}, error) {
	var rstData map[string]interface{}
	err := json.Unmarshal([]byte(rst), &rstData)
	return rstData, err
}

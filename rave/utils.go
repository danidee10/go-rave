// utils and helper functions

package rave

import (
	"encoding/json"
)

func mapToJSON(mapData map[string]interface{}) []byte {
	jsonBytes, err := json.Marshal(mapData)
	if err != nil {
		panic(err)
	}

	return jsonBytes
}

func jsonToMap(jsonData []byte) map[string]interface{} {
	jsonMap := make(map[string]interface{})

	err := json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		panic(err)
	}

	return jsonMap
}

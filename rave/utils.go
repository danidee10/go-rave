// utils and helper functions

package rave

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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

func jsonToInterfaceList(jsonData []byte) []interface{} {
	var returnData []interface{}
	err := json.Unmarshal(jsonData, &returnData)
	if err != nil {
		panic(err)
	}

	return returnData
}

// Check if an array of keys is set in map
func checkRequiredParameters(params map[string]interface{}, keys []string) {
	for _, key := range keys {
		if _, ok := params[key]; !ok {
			log.Fatalf("%s is a required parameter for this method", key)
		}
	}
}

// MakePostRequest : make s post request with the Content-Type set to application/json
func MakePostRequest(URL string, data map[string]interface{}) map[string]interface{} {
	postData := mapToJSON(data)
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(postData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return jsonToMap(body)
}

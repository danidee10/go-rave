// utils and helper functions

package rave

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/antonholmquist/jason"
)

// MapToJSON : Convert map[string]interface{} to JSON
func mapToJSON(mapData map[string]interface{}) []byte {
	jsonBytes, err := json.Marshal(mapData)
	if err != nil {
		panic(err)
	}

	return jsonBytes
}

// Check if an array of keys is set in map
func checkRequiredParameters(params map[string]interface{}, keys []string) error {
	for _, key := range keys {
		if _, ok := params[key]; !ok {
			return fmt.Errorf("\"%s\" is a required parameter for this method", key)
		}
	}

	return nil
}

// MakePostRequest : make s post request with the Content-Type set to application/json
func MakePostRequest(URL string, data map[string]interface{}) ([]byte, error) {
	postData := mapToJSON(data)
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(postData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = handleAPIErrors(resp, body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// handle errors raised by the API's, this include's non 200 Errors
// and Errors for missing or invalid parameters
func handleAPIErrors(response *http.Response, body []byte) error {
	v, _ := jason.NewObjectFromBytes(body)
	status, _ := v.GetString("status")

	if status != "success" {
		errorMessage, _ := v.GetString("message")
		return fmt.Errorf("%s Status Code: %d", errorMessage, response.StatusCode)
	}

	return nil
}

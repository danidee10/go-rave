/* This file contains functions/methods that are not related to payments directly */

package rave

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

// CalculateIntegrityCheckSum : Calculates the integrity checksum of the data required by the browser
func (r Rave) CalculateIntegrityCheckSum(data map[string]interface{}) string {
	// sort the map
	sortedKeys := []string{}
	sortedValues := []string{}

	for key := range data {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		value := data[key]
		// convert all values to strings before appending
		sortedValues = append(sortedValues, fmt.Sprint(value))
	}

	// convert sortedValues to string slice so we can join
	for key, value := range sortedValues {
		sortedValues[key] = value
	}

	// concatenate the sorted values
	sha256Payload := strings.Join(sortedValues[:], "")

	// join with secret key
	sha256Payload += r.GetSecretKey()

	// Generate a sha256 hash and convert the bytes to hex
	integrityCheckSum := fmt.Sprintf("%x", sha256.Sum256([]byte(sha256Payload)))

	return integrityCheckSum
}

// GetFees : Get fees to be charged for a particular amount/currency
func (r Rave) GetFees(data map[string]interface{}) ([]byte, error) {
	err := checkRequiredParameters(data, []string{"amount", "currency"})

	data["PBFPubKey"] = r.GetPublicKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/fee"

	response, err := MakePostRequest(URL, data)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ListBanks : List Nigerian banks.
func (r Rave) ListBanks() ([]byte, error) {
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/flwpbf-banks.js?json=1"
	response, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	return body, nil
}

package rave

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"
)

// Rave : Base Rave type
type rave struct {
	publicKey string
	secretKey string
	Live      bool
	LiveURL   string
	TestURL   string
}

// getBaseURL : Returns the Correct URL based on Live status.
func (r rave) getBaseURL() string {
	if r.Live {
		return r.LiveURL
	}

	return r.TestURL
}

// GetPublicKey: Get Rave Public key
func (r rave) GetPublicKey() string {
	publicKey, found := os.LookupEnv("RAVE_PUBLICKEY")
	if !found {
		log.Fatal("You must set the \"RAVE_PUBLICKEY\" environment variable")
	}

	return publicKey
}

// GetSecretKey: Get Rave Secret key
func (r rave) GetSecretKey() string {
	secKey, found := os.LookupEnv("RAVE_SECKEY")
	if !found {
		log.Fatal("You must set the \"RAVE_SECKEY\" environment variable")
	}

	return secKey
}

// ChargeCard : Sends a Card request and determines the validation flow to be used
func (r rave) ChargeCard(chargeData map[string]interface{}) map[string]interface{} {
	postData := r.setUpCharge(chargeData)
	response := r.chargeCard(postData)

	// If suggested_auth == "PIN" was returned in the response
	// Encrypt the client's data with the otp details and make another request
	suggestedAuthData := response["data"]
	if reflect.DeepEqual(suggestedAuthData, map[string]interface{}{"suggested_auth": "PIN"}) {
		chargeData["suggested_auth"] = "PIN"
		checkRequiredParameters(chargeData, []string{"pin"})

		postData := r.setUpCharge(chargeData)
		response := r.chargeCard(postData)

		return response
	}

	return response
}

// Encrypts and setup a charge (Payment/account) with the secret key and algorithm
func (r rave) setUpCharge(chargeData map[string]interface{}) map[string]interface{} {
	chargeJSON := mapToJSON(chargeData)
	encryptedchargeData := r.Encrypt3Des(string(chargeJSON[:]))

	data := map[string]interface{}{
		"PBFPubKey": r.GetPublicKey(),
		"client":    encryptedchargeData,
		"alg":       "3DES-24",
	}

	return data
}

// ChargeCard: Contains the actual logic for making requests to the charge endpoint
func (r rave) chargeCard(postData map[string]interface{}) map[string]interface{} {
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/charge"

	response := MakePostRequest(URL, postData)

	return response
}

// ValidateCharge : Validate a card charge using OTP
func (r rave) ValidateCharge(data map[string]interface{}) map[string]interface{} {

	data["PBFPubKey"] = r.GetPublicKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/validatecharge"

	response := MakePostRequest(URL, data)

	return response
}

// ValidateAccountCharge : Validate an account charge using OTP
func (r rave) ValidateAccountCharge(data map[string]interface{}) map[string]interface{} {

	data["PBFPubKey"] = r.GetPublicKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/validate"

	response := MakePostRequest(URL, data)

	return response
}

// VerifyTransaction: Verify a transaction using "flw_ref" or "tx_ref"
func (r rave) VerifyTransaction(data map[string]interface{}) map[string]interface{} {
	data["SECKEY"] = r.GetSecretKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/verify"

	response := MakePostRequest(URL, data)

	return response
}

func (r rave) XrequeryTransactionVerification(data map[string]interface{}) map[string]interface{} {
	data["SECKEY"] = r.GetSecretKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/xrequery"

	response := MakePostRequest(URL, data)

	return response
}

// PreauthorizeCard : This is just a wrapper arond the Charge function that automatically
// sets "charge_type" to "pre_auth"
func (r rave) PreauthorizeCard(chargeData map[string]interface{}) map[string]interface{} {
	chargeData["charge_type"] = "pre_auth"

	ret := r.ChargeCard(chargeData)

	return ret

}

func (r rave) Capture(data map[string]interface{}) map[string]interface{} {
	data["SECKEY"] = r.GetSecretKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/capture"

	response := MakePostRequest(URL, data)

	return response
}

// RefundOrVoid : Refund or void a captured amount
func (r rave) RefundOrVoid(data map[string]interface{}) map[string]interface{} {
	data["SECKEY"] = r.GetSecretKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/refundorvoid"

	response := MakePostRequest(URL, data)

	return response
}

// GetFees
func (r rave) GetFees(data map[string]interface{}) map[string]interface{} {
	data["PBFPubKey"] = r.GetPublicKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/fee"

	response := MakePostRequest(URL, data)

	return response
}

// Refund : Refund direct charges
func (r rave) Refund(data map[string]interface{}) map[string]interface{} {
	data["seckey"] = r.GetSecretKey()
	URL := r.getBaseURL() + "/gpx/merchant/transactions/refund"

	response := MakePostRequest(URL, data)

	return response
}

// ListBanks : List Nigerian banks.
func (r rave) ListBanks() []interface{} {
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/flwpbf-banks.js?json=1"
	response, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	returnData := jsonToInterfaceList(body)

	return returnData
}

// CalculateIntegrityCheckSum : Calculates the integrity checksum of the data required by the browser
func (r rave) CalculateIntegrityCheckSum(data map[string]interface{}) string {
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

// NewRave : Constructor for rave struct
func NewRave() rave {
	Rave := rave{}
	Rave.TestURL = "http://flw-pms-dev.eu-west-1.elasticbeanstalk.com"
	Rave.LiveURL = "https://api.ravepay.co"

	// default mode is development
	Rave.Live = false

	return Rave
}

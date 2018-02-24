package rave

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
)

// Rave : Base Rave type
type rave struct {
	publicKey string
	secretKey string
	Live      bool
	LiveURL   string
	TestURL   string
	Prefix    string // Your application name e.g FLUTTERWAVE
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
	key := r.GetKey(r.GetSecretKey())
	chargeJSON := mapToJSON(chargeData)
	encryptedchargeData := r.Encrypt3Des(key, string(chargeJSON[:]))

	data := map[string]interface{}{
		"PBFPubKey": r.GetPublicKey(),
		"client":    encryptedchargeData,
		"alg":       "3DES-24",
	}

	return data
}

// ChargeCard: Contains the actual logic for making requests to the charge endpoint
func (r rave) chargeCard(postData map[string]interface{}) map[string]interface{} {
	URL := r.getBaseURL() + "/getpaidx/api/charge"

	response := MakePostRequest(URL, postData)

	return response
}

// ValidateCharge : Validate a card charge using OTP
func (r rave) ValidateCharge(data map[string]interface{}) map[string]interface{} {

	data["PBFPubKey"] = r.GetPublicKey()
	URL := r.getBaseURL() + "/getpaidx/api/validatecharge"

	response := MakePostRequest(URL, data)

	return response
}

// ValidateAccountCharge : Validate an account charge using OTP
func (r rave) ValidateAccountCharge(data map[string]interface{}) map[string]interface{} {

	data["PBFPubKey"] = r.GetPublicKey()
	URL := r.getBaseURL() + "/getpaidx/api/validate"

	response := MakePostRequest(URL, data)

	return response
}

// VerifyTransaction: Verify a transaction using "flw_ref" or "tx_ref"
func (r rave) VerifyTransaction(data map[string]interface{}) map[string]interface{} {
	data["SECKEY"] = r.GetSecretKey()
	URL := r.getBaseURL() + "/getpaidx/api/verify"

	response := MakePostRequest(URL, data)

	return response
}

func (r rave) XrequeryTransactionVerification(data map[string]interface{}) map[string]interface{} {
	data["SECKEY"] = r.GetSecretKey()
	URL := r.getBaseURL() + "/getpaidx/api/xrequery"

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
	URL := r.getBaseURL() + "/getpaidx/api/capture"

	response := MakePostRequest(URL, data)

	return response
}

// RefundOrVoid : Refund or void a captured amount
func (r rave) RefundOrVoid(data map[string]interface{}) map[string]interface{} {
	data["SECKEY"] = r.GetSecretKey()
	URL := r.getBaseURL() + "/getpaidx/api/refundorvoid"

	response := MakePostRequest(URL, data)

	return response
}

// GetFees
func (r rave) GetFees(data map[string]interface{}) map[string]interface{} {
	data["PBFPubKey"] = r.GetPublicKey()
	URL := r.getBaseURL() + "/getpaidx/api/fee"

	response := MakePostRequest(URL, data)

	return response
}

// Refund : Refund direct charges
func (r rave) Refund(data map[string]interface{}) map[string]interface{} {
	data["seckey"] = r.GetSecretKey()
	var URL string
	if r.Live {
		URL = "http://flw-pms-dev.eu-west-1.elasticbeanstalk.com/gpx/merchant/transactions/refund"
	} else {
		URL = "http://flw-pms-dev.eu-west-1.elasticbeanstalk.com/gpx/merchant/transactions/refund"
	}

	response := MakePostRequest(URL, data)

	return response
}

// ListBanks : List Nigerian banks.
func (r rave) ListBanks() []interface{} {
	URL := r.getBaseURL() + "/getpaidx/api/flwpbf-banks.js?json=1"
	response, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	returnData := jsonToInterfaceList(body)

	return returnData
}

// NewRave : Constructor for rave struct
func NewRave() rave {
	Rave := rave{}
	Rave.TestURL = "http://flw-pms-dev.eu-west-1.elasticbeanstalk.com/flwv3-pug"

	return Rave
}

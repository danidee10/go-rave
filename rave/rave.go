package rave

import (
	"bytes"
	"io/ioutil"
	"net/http"
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
	return "FLWPUBK-519ac5f00bd2855a2f25c556c01888cd-X"
}

// GetSecretKey: Get Rave Secret key
func (r rave) GetSecretKey() string {
	return "FLWSECK-dfc21e7469c575750846ee73820b6374-X"
}

func (r rave) AccountChargeNigeria() {

}

func (r rave) AccountChargeInternational() {

}

// ChargeCard: Charge a card
func (r rave) ChargeCard(clientData map[string]interface{}) map[string]interface{} {

	clientJSON := mapToJSON(clientData)

	key := r.GetKey(r.GetSecretKey())
	encryptedClientData := r.Encrypt3Des(key, string(clientJSON[:]))

	data := map[string]interface{}{
		"PBFPubKey": r.GetPublicKey(),
		"client":    encryptedClientData,
		"alg":       "3DES-24",
	}

	postData := mapToJSON(data)

	URL := r.getBaseURL() + "/getpaidx/api/charge"
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

// ValidateCharge : Validate a card charge using OTP
func (r rave) ValidateCharge(data map[string]interface{}) map[string]interface{} {

	data["PBFPubKey"] = r.GetPublicKey()
	postData := mapToJSON(data)

	URL := r.getBaseURL() + "/getpaidx/api/validatecharge"
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

// ValidateAccountCharge : Validate an account charge using OTP
func (r rave) ValidateAccountCharge(data map[string]interface{}) map[string]interface{} {

	data["PBFPubKey"] = r.GetPublicKey()
	postData := mapToJSON(data)

	URL := r.getBaseURL() + "/getpaidx/api/validate"
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

// VerifyTransaction: Verify a transaction using "flw_ref" or "tx_ref"
func (r rave) VerifyTransaction(data map[string]interface{}) map[string]interface{} {
	data["SECKEY"] = r.GetSecretKey()

	postData := mapToJSON(data)

	URL := r.getBaseURL() + "/getpaidx/api/verify"
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

func (r rave) XrequeryTransactionVerification(data map[string]interface{}) map[string]interface{} {
	data["SECKEY"] = r.GetSecretKey()

	postData := mapToJSON(data)

	URL := r.getBaseURL() + "/getpaidx/api/xrequery"

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

func (r rave) RetryTransaction() {

}

// PreauthorizeCard : This is just a wrapper arond the Charge function that automatically
// sets "charge_type" to "pre_auth"
func (r rave) PreauthorizeCard(clientData map[string]interface{}) map[string]interface{} {
	clientData["charge_type"] = "pre_auth"

	ret := r.ChargeCard(clientData)

	return ret

}

func (r rave) Capture(data map[string]interface{}) map[string]interface{} {
	data["SECKEY"] = r.GetSecretKey()

	postData := mapToJSON(data)

	URL := r.getBaseURL() + "/getpaidx/api/capture"

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

// RefundOrVoid : Refund or void a captured amount
func (r rave) RefundOrVoid(data map[string]interface{}) map[string]interface{} {
	data["SECKEY"] = r.GetSecretKey()

	postData := mapToJSON(data)

	URL := r.getBaseURL() + "/getpaidx/api/refundorvoid"

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

// GetFees
func (r rave) GetFees(data map[string]interface{}) map[string]interface{} {
	data["PBFPubKey"] = r.GetPublicKey()

	postData := mapToJSON(data)

	URL := r.getBaseURL() + "/getpaidx/api/fee"
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

// Refund : Refund direct charges
func (r rave) Refund(data map[string]interface{}) map[string]interface{} {
	data["seckey"] = r.GetSecretKey()

	postData := mapToJSON(data)

	var URL string
	if r.Live {
		URL = "http://flw-pms-dev.eu-west-1.elasticbeanstalk.com/gpx/merchant/transactions/refund"
	} else {
		URL = "http://flw-pms-dev.eu-west-1.elasticbeanstalk.com/gpx/merchant/transactions/refund"
	}
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

// ListBanks : List Nigerian banks.
func (r rave) ListBanks() *http.Response {
	URL := r.getBaseURL() + "/getpaidx/api/flwpbf-banks.js?json=1"
	response, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	return response
}

// NewRave : Constructor for rave struct
func NewRave() rave {
	Rave := rave{}
	Rave.TestURL = "http://flw-pms-dev.eu-west-1.elasticbeanstalk.com/flwv3-pug"

	return Rave
}

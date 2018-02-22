package rave

import "net/http"

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
	return ""
}

// GetSecretKey: Get Rave Secret key
func (r rave) GetSecretKey() string {
	return "FLWSECK-6b32914d4d60c10d0ef72bdad734134a-X"
}

func (r rave) AccountChargeNigeria() {

}

func (r rave) AccountChargeInternational() {

}

func (r rave) CardCharge() {

}
func (r rave) VerifyTransaction() {

}

func (r rave) RetryTransaction() {

}

func (r rave) PreauthCaptureRefund() {

}

func (r rave) PayWithUSSDMCash() {

}

// ListBanks : List Nigerian banks.
func (r rave) ListBanks() *http.Response {
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/flwpbf-banks.js?json=1"
	response, err := http.Get(URL)
	if err != nil {
		panic(err)
	}

	return response
}

func (r rave) GetFees() {

}

func (r rave) CalculateIntegrityCheckSum() {

}

// NewRave : Constructor for rave struct
func NewRave() rave {
	Rave := rave{}
	Rave.TestURL = "https://rave-api-v2.herokuapp.com"

	return Rave
}

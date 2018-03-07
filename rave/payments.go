/* This file contains methods/functions that deal with payments (Charges) */

package rave

import (
	"github.com/antonholmquist/jason"
)

// ChargeCard : Sends a Card request and determine the validation flow to be used
func (r Rave) ChargeCard(chargeData map[string]interface{}) ([]byte, error) {
	err := checkRequiredParameters(chargeData, []string{
		"cardno", "cvv", "expirymonth", "expiryyear", "amount", "email",
		"phonenumber", "firstname", "lastname", "IP", "txRef", "redirect_url",
	})
	if err != nil {
		return nil, err
	}

	postData := r.setUpCharge(chargeData)
	response, err := r.charge(postData)
	if err != nil {
		return nil, err
	}

	// If suggested_auth == "PIN" was returned in the response
	// Encrypt the client's data with the otp details and make another request
	suggestedAuthData, _ := jason.NewObjectFromBytes(response)
	data, _ := suggestedAuthData.GetObject("data")
	suggestedAuth, _ := data.GetString("suggested_auth")

	if suggestedAuth == "PIN" {
		chargeData["suggested_auth"] = "PIN"
		err := checkRequiredParameters(chargeData, []string{"pin"})
		if err != nil {
			return nil, err
		}

		response, err = r.ChargeCard(chargeData)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}

// Encrypts and setup a charge (Payment/account) with the secret key and algorithm
func (r Rave) setUpCharge(chargeData map[string]interface{}) map[string]interface{} {
	chargeJSON := mapToJSON(chargeData)
	encryptedchargeData := r.Encrypt3Des(string(chargeJSON[:]))

	data := map[string]interface{}{
		"PBFPubKey": r.GetPublicKey(),
		"client":    encryptedchargeData,
		"alg":       "3DES-24",
	}

	return data
}

// charge: Contains the actual logic for making requests to the charge endpoint
func (r Rave) charge(data map[string]interface{}) ([]byte, error) {
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/charge"

	response, err := MakePostRequest(URL, data)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ValidateCharge : Validate a card charge using OTP
func (r Rave) ValidateCharge(data map[string]interface{}) ([]byte, error) {
	err := checkRequiredParameters(data, []string{"transaction_reference", "otp"})
	if err != nil {
		return nil, err
	}

	data["PBFPubKey"] = r.GetPublicKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/validatecharge"

	response, err := MakePostRequest(URL, data)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ChargeAccount : Charge a Local (Nigerian) or South African Bank Account
func (r Rave) ChargeAccount(data map[string]interface{}) ([]byte, error) {
	err := checkRequiredParameters(data, []string{
		"accountnumber", "accountbank", "email", "phonenumber",
		"firstname", "lastname", "IP", "txRef", "payment_type",
	})
	if err != nil {
		return nil, err
	}

	postData := r.setUpCharge(data)
	response, err := r.charge(postData)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ValidateAccountCharge : Validate an account charge using OTP
func (r Rave) ValidateAccountCharge(data map[string]interface{}) ([]byte, error) {
	err := checkRequiredParameters(data, []string{"transactionreference", "otp"})
	if err != nil {
		return nil, err
	}

	data["PBFPubKey"] = r.GetPublicKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/validate"

	response, err := MakePostRequest(URL, data)
	if err != nil {
		return nil, err
	}

	return response, nil
}

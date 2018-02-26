// Tests for the rave package

package rave

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/antonholmquist/jason"
)

//=============================================================================
// Test Setup

var Rave rave

// Setup test suite
func TestMain(m *testing.M) {
	Rave = NewRave()
	Rave.Live = false
	fmt.Println("Running tests...")

	os.Exit(m.Run())
}

// End test setup
// ============================================================================

// Test the encryption function
func TestEncryption(t *testing.T) {
	t.Parallel()

	assertEqual(t, Rave.Encrypt3Des("Hello world"), "fus4LnqrvKWXqm7wueoj2Q==")
}

func TestSuggestedAuthPin(t *testing.T) {
	t.Parallel()

	masterCard := map[string]interface{}{
		"name": "suggestedAuthPin", "cardno": "5438898014560229", "currency": "NGN",
		"country": "NG", "cvv": "789", "amount": "300", "expiryyear": "19",
		"expirymonth": "09", "pin": "3310", "email": "suggestedAuthPin@flutter.co",
		"IP": "103.238.105.185", "txRef": "MXX-ASC-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
		"redirect_url":       "http://127.0.0.1",
	}

	response, _ := Rave.ChargeCard(masterCard)

	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	authModelUsed, _ := data.GetString("authModelUsed")

	assertEqual(t, authModelUsed, "PIN")
}

// It should raise an error if the pin wasn't passed and the suggested_auth is "PIN"
func TestSuggestedAuthPinRaisesError(t *testing.T) {
	t.Parallel()

	masterCard := map[string]interface{}{
		"name": "suggestedAuth", "cardno": "5438898014560229", "currency": "NGN",
		"country": "NG", "cvv": "789", "amount": "300", "expiryyear": "19",
		"expirymonth": "09", "email": "suggestedAuthr@flutter.co", "IP": "103.238.105.185",
		"txRef": "MXX-ASC-4578", "device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
		"redirect_url": "http://127.0.0.1",
	}
	_, err := Rave.ChargeCard(masterCard)

	assertEqual(t, err.Error(), "\"pin\" is a required parameter for this method")
}

// Method should return "VBVSECURECODE" or "AVS_VBVSECURECODE" as the suggestedAuth
func TestSuggestedAuth3DesSecurePayment(t *testing.T) {
	t.Parallel()

	visaCard := map[string]interface{}{
		"name": "Suggested3DesSecurePayment", "cardno": "4556052704172643", "currency": "USD",
		"country": "US", "cvv": "899", "amount": "1000", "expiryyear": "19",
		"expirymonth": "09", "email": "Suggested3DesSecurePayment@flutter.co", "txRef": "TXT",
		"redirect_url": "http://127.0.0.1",
	}

	response, _ := Rave.ChargeCard(visaCard)

	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	authModelUsed, _ := data.GetString("authModelUsed")

	if authModelUsed != "AVS_VBVSECURECODE" && authModelUsed != "VBVSECURECODE" {
		t.Errorf("The authModelUsed is not 'AVS_VBVSECURECODE' or 'VBVSECURECODE', it's %s", authModelUsed)
	}
}

// Should raise an error because of the missing redirect_url parameter
func TestSuggestedAuth3DesSecurePaymentRaisesError(t *testing.T) {
	t.Parallel()

	visaCard := map[string]interface{}{
		"name": "Suggested3DesSecurePaymentRaisesError", "cardno": "4556052704172643",
		"currency": "USD", "country": "US", "cvv": "899", "amount": "1000", "expiryyear": "19",
		"expirymonth": "09", "email": "Suggested3DesSecurePaymentRaisesError@flutter.co",
		"txRef": "TXT",
	}

	_, err := Rave.ChargeCard(visaCard)

	assertEqual(t, err.Error(), "\"redirect_url\" is a required parameter for this method")
}

func TestMasterCardPaymentWithPin(t *testing.T) {
	t.Parallel()

	masterCard := map[string]interface{}{
		"name": "paymentWithPin", "cardno": "5438898014560229", "currency": "NGN",
		"country": "NG", "cvv": "789", "amount": "300", "expiryyear": "19",
		"expirymonth": "09", "suggested_auth": "pin", "pin": "3310",
		"email": "paymentWithPin@flutter.co", "IP": "103.238.105.185", "txRef": "MXX-ASC-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
		"redirect_url":       "http://127.0.0.1",
	}

	response, _ := Rave.ChargeCard(masterCard)

	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	authModelUsed, _ := data.GetString("authModelUsed")

	assertEqual(t, authModelUsed, "PIN")
}

func testVerveCardPaymentWithPin(t *testing.T) {
	t.Parallel()

	verveCard := map[string]interface{}{
		"name": "verve", "cardno": "5061020000000000094", "currency": "NGN",
		"country": "NG", "cvv": "347", "amount": "300", "expiryyear": "20",
		"expirymonth": "07", "suggested_auth": "pin", "pin": "1111",
		"email": "verve@flutter.co", "IP": "103.238.105.185", "txRef": "MXX-ASC-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
		"redirect_url":       "http://127.0.0.1",
	}

	response, err := Rave.ChargeCard(verveCard)
	if err != nil {
		fmt.Println(err.Error())
	}

	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	authModelUsed, _ := data.GetString("authModelUsed")

	assertEqual(t, authModelUsed, "PIN")

	fmt.Println(string(response[:]))
}

func TestVisaPaymentWith3DSecure(t *testing.T) {
	t.Parallel()

	visaCard := map[string]interface{}{
		"name": "visa", "cardno": "4187427415564246", "currency": "NGN",
		"country": "NG", "cvv": "828", "amount": "300", "expiryyear": "19",
		"expirymonth": "09", "email": "visa@flutter.co", "IP": "103.238.105.185",
		"txRef": "MXX-ASC-4578", "device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
		"redirect_url": "http://127.0.0.1",
	}
	response, _ := Rave.ChargeCard(visaCard)

	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	authURL, _ := data.GetString("authurl")

	if authURL == "" {
		t.Error("authurl not found in the response")
	}
}

// Every method (That makes use of the MakePostRequest method)
// should return a response (as map[string]interface) for any failed request
// When the API returns an error
func TestErrorResponse(t *testing.T) {
	t.Parallel()

	// Make a request without including the cvv
	verveCard := map[string]interface{}{
		"name": "cvvError", "cardno": "5061020000000000094", "currency": "NGN",
		"country": "NG", "amount": "300", "expiryyear": "20",
		"expirymonth": "07", "suggested_auth": "pin", "pin": "1111",
		"email": "cvvError@flutter.co", "IP": "103.238.105.185", "txRef": "MXX-ASC-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
		"redirect_url":       "http://127.0.0.1",
	}
	_, err := Rave.ChargeCard(verveCard)

	if err == nil {
		t.Error("'TestErrorResponse' didn't raise an error")
	}

	errorString := err.Error()
	if errorString != "cvv is required Status Code: 400" {
		t.Errorf("Method didn't raise 'cvv is required error' instead it raised %s", errorString)
	}
}

// We should get a list of all Nigerian banks we can charge
func TestListBanks(t *testing.T) {
	t.Parallel()

	response, _ := Rave.ListBanks()

	var banks []map[string]string
	json.Unmarshal(response, &banks)

	// Ensure that access bank is in the response
	accessBank := banks[0]
	if accessBank["bankname"] != "ACCESS BANK NIGERIA" || accessBank["bankcode"] != "044" {
		t.Error("Access Bank wasn't listed")
		fmt.Println(string(response[:]))
	}
}

// Test that a charge on a card can be validated using OTP
func TestChargeCard(t *testing.T) {
	t.Parallel()

	// Initialize the transaction and get a valid transaction reference
	masterCard := map[string]interface{}{
		"name": "chargeCard", "cardno": "5438898014560229", "currency": "NGN",
		"country": "NG", "cvv": "789", "amount": "300", "expiryyear": "19",
		"expirymonth": "09", "suggested_auth": "pin", "pin": "3310",
		"email": "chargeCard@flutter.co", "IP": "103.238.105.185", "txRef": "MXX-ASC-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
		"redirect_url":       "http://127.0.0.1",
	}

	response, _ := Rave.ChargeCard(masterCard)

	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	transactionReference, _ := data.GetString("flwRef")

	transaction := map[string]interface{}{
		"transaction_reference": transactionReference,
		"otp": "12345",
	}

	response, _ = Rave.ValidateCharge(transaction)
	v, _ = jason.NewObjectFromBytes(response)
	successMessage, _ := v.GetString("message")
	data, _ = v.GetObject("data")
	tx, _ := data.GetObject("tx")
	chargedAmount, _ := tx.GetInt64("charged_amount")

	if successMessage != "Charge Complete" || chargedAmount != 300 {
		t.Error("Card Charge failed")
		fmt.Println(successMessage, chargedAmount)
	}
}

// Verify the status of a transaction
func TestVerifyTransaction(t *testing.T) {
	t.Parallel()

	// Initialize the transaction and get a valid transaction reference
	masterCard := map[string]interface{}{
		"name": "verifyTransaction", "cardno": "5438898014560229", "currency": "NGN",
		"country": "NG", "cvv": "789", "amount": "300", "expiryyear": "19",
		"expirymonth": "09", "suggested_auth": "pin", "pin": "3310",
		"email": "verifyTransaction@flutter.co", "IP": "103.238.105.185", "txRef": "MXX-ASC-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
		"redirect_url":       "http://127.0.0.1",
	}

	response, _ := Rave.ChargeCard(masterCard)

	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	transactionReference, _ := data.GetString("flwRef")
	currency, _ := data.GetString("currency")

	// Pay for the transaction
	transaction := map[string]interface{}{
		"transaction_reference": transactionReference,
		"otp": "12345",
	}
	response, _ = Rave.ValidateCharge(transaction)

	// Verify the transaction
	transaction = map[string]interface{}{
		"flw_ref": transactionReference, "normalize": "1",
		"currency": currency, "amount": "1000",
	}
	response, err := Rave.VerifyTransaction(transaction)
	if err != nil {
		t.Error(err.Error())
	}
}

// Verify the status of a transaction using XRequery
func TestXrequeryTransactionVerification(t *testing.T) {
	// Initialize the transaction and get a valid transaction reference
	masterCard := map[string]interface{}{
		"name": "xrequery", "cardno": "5438898014560229", "currency": "NGN",
		"country": "NG", "cvv": "789", "amount": "5300", "expiryyear": "19",
		"expirymonth": "09", "suggested_auth": "pin", "pin": "3310",
		"email": "xrequery@flutter.co", "IP": "103.238.105.111", "txRef": "abcdef",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe031234e",
		"redirect_url":       "http://127.0.0.1",
	}

	response, _ := Rave.ChargeCard(masterCard)

	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	transactionReference, _ := data.GetString("flwRef")
	currency, _ := data.GetString("currency")

	// Pay for the transaction
	transaction := map[string]interface{}{
		"transaction_reference": transactionReference,
		"otp": "12345",
	}
	response, _ = Rave.ValidateCharge(transaction)

	// Verify the transaction
	// flw_ref is needed for verification
	transaction = map[string]interface{}{
		"flw_ref": transactionReference, "tx_ref": "abcdef",
		"last_attempt": "1", "only_attempt": "1",
		"currency": currency, "amount": "5300",
	}
	response, err := Rave.XrequeryTransactionVerification(transaction)
	if err != nil {
		t.Error(err.Error())
	}
}

// Test Get fees endpoint
func testGetFees(t *testing.T) {
	t.Parallel()

	data := map[string]interface{}{
		"amount": "5300", "currency": "NGN",
	}

	response, _ := Rave.GetFees(data)

	fmt.Println(string(response[:]))
}

// Test if the CalculateIntegrityCheckSum function returns valid results
func TestCalculateIntegrityCheckSum(t *testing.T) {
	t.Parallel()

	data := map[string]interface{}{
		"PBFPubKey":          "FLWPUBK-e634d14d9ded04eaf05d5b63a0a06d2f-X",
		"amount":             20,
		"payment_method":     "both",
		"custom_description": "Pay Internet",
		"custom_logo":        "http://localhost/payporte-3/skin/frontend/ultimo/shoppy/custom/images/logo.svg",
		"custom_title":       "Shoppy Global systems",
		"country":            "NG",
		"currency":           "NGN",
		"customer_email":     "user@example.com",
		"customer_firstname": "Temi",
		"customer_lastname":  "Adelewa",
		"customer_phone":     "234099940409",
		"txref":              "MG-1500041286295",
	}

	// set Rave seckey environment variable so it matches the expected result
	oldSecKey, found := os.LookupEnv("RAVE_SECKEY")
	if !found {
		log.Fatal("You must set the \"RAVE_SECKEY\" environment variable")
	}

	err := os.Setenv("RAVE_SECKEY", "FLWSECK-bb971402072265fb156e90a3578fe5e6-X")
	if err != nil {
		panic(err)
	}

	integrityChecksum := Rave.CalculateIntegrityCheckSum(data)

	assertEqual(t, integrityChecksum, "a14ac4eba0902e8fd6b5fdf542f46d6efc18885a63c3d5f100c26715c7c8d8f4")

	// set "RAVE_SECKEY" to it's old value
	os.Setenv("RAVE_SECKEY", oldSecKey)
}

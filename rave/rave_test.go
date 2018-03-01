// Tests for the rave package

package rave

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/antonholmquist/jason"
)

//=============================================================================
// Test Setup

var rave Rave

var DefaultPublicKey = rave.GetPublicKey()
var DefaultSecretKey = rave.GetSecretKey()

// Setup test suite
func TestMain(m *testing.M) {
	rave = NewRave()
	rave.Live = false
	fmt.Println("Running tests...")

	os.Exit(m.Run())
}

// End test setup
// ============================================================================

// Test the encryption function
func TestEncryption(t *testing.T) {
	t.Parallel()

	assertEqual(t, rave.Encrypt3Des("Hello world"), "fus4LnqrvKWXqm7wueoj2Q==")
}

func TestSuggestedAuthPin(t *testing.T) {
	t.Parallel()

	masterCard := map[string]interface{}{
		"name": "suggestedAuthPin", "cardno": "5438898014560229", "currency": "NGN",
		"country": "NG", "cvv": "789", "amount": "300", "expiryyear": "19",
		"expirymonth": "09", "pin": "3310", "email": "TestSuggestedAuth@flutter.co",
		"firstname": "suggested", "lastname": "auth", "phonenumber": "081245554343",
		"IP": "103.238.105.185", "txRef": "MXX-AYT-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
		"redirect_url":       "http://127.0.0.1",
	}

	response, _ := rave.ChargeCard(masterCard)

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
		"expirymonth": "09", "email": "TestSuggestedAuthPinRaisesError@flutter.co",
		"firstname": "suggested_pin", "lastname": "raises_error", "phonenumber": "081245554343",
		"IP": "103.238.105.185", "txRef": "MXX-AYT-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c", "redirect_url": "http://127.0.0.1",
	}
	_, err := rave.ChargeCard(masterCard)

	assertEqual(t, err.Error(), "\"pin\" is a required parameter for \"ChargeCard\"")
}

// Method should return "VBVSECURECODE" or "AVS_VBVSECURECODE" as the suggestedAuth
func TestSuggestedAuth3DesSecurePayment(t *testing.T) {
	t.Parallel()

	visaCard := map[string]interface{}{
		"name": "Suggested3DesSecurePayment", "cardno": "4556052704172643", "currency": "USD",
		"country": "US", "cvv": "899", "amount": "1000", "expiryyear": "19",
		"expirymonth": "09", "email": "TestSuggestedAuth3Des@flutter.co",
		"firstname": "suggested_auth", "lastname": "3Desauth", "phonenumber": "081245554343",
		"IP": "103.238.105.185", "txRef": "MXX-AYT-4578",
		"redirect_url": "http://127.0.0.1",
	}

	response, _ := rave.ChargeCard(visaCard)

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
		"expirymonth": "09", "email": "TestSuggestedAuth3DesRaisesError@flutter.co",
		"firstname": "suggested3DesRaisesError", "lastname": "3DesRaisesError",
		"phonenumber": "081245554343", "IP": "103.238.105.185", "txRef": "MXX-AYT-4578",
	}

	_, err := rave.ChargeCard(visaCard)

	assertEqual(t, err.Error(), "\"redirect_url\" is a required parameter for this method")
}

func TestMasterCardPaymentWithPin(t *testing.T) {
	t.Parallel()

	masterCard := map[string]interface{}{
		"name": "paymentWithPin", "cardno": "5438898014560229", "currency": "NGN",
		"country": "NG", "cvv": "789", "amount": "300", "expiryyear": "19",
		"expirymonth": "09", "suggested_auth": "pin", "pin": "3310",
		"email": "PaymentWithPin@flutter.co", "firstname": "payment",
		"lastname": "with_pin", "phonenumber": "081245554343",
		"IP": "103.238.105.185", "txRef": "MXX-AYT-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
		"redirect_url":       "http://127.0.0.1",
	}

	response, _ := rave.ChargeCard(masterCard)

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
		"email": "verve@flutter.co", "firstname": "verve", "lastname": "verve",
		"phonenumber": "081245554343", "IP": "103.238.105.185", "txRef": "MXX-AYT-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
		"redirect_url":       "http://127.0.0.1",
	}

	response, err := rave.ChargeCard(verveCard)
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
		"expirymonth": "09", "email": "visa@flutter.co",
		"firstname": "visa", "lastname": "visa", "phonenumber": "081245554343",
		"IP": "103.238.105.185", "txRef": "MXX-AYT-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
		"redirect_url":       "http://127.0.0.1",
	}
	response, _ := rave.ChargeCard(visaCard)

	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	authURL, _ := data.GetString("authurl")

	if authURL == "" {
		t.Error("authurl not found in the response")
	}
}

// Test AccountCharge
func TestAccountCharge(t *testing.T) {
	t.Parallel()

	// get access bank details
	response, _ := rave.ListBanks()

	var banks []map[string]string
	json.Unmarshal(response, &banks)

	accessBankCode := banks[0]["bankcode"]

	accountDetails := map[string]interface{}{
		"accountnumber": "0690000031", "accountbank": accessBankCode, "currency": "NGN",
		"country": "NG", "amount": 5000, "email": "TestAccessBankCharge@gmail.com",
		"phonenumber": "08123456787", "firstname": "Access", "lastname": "AccessBank",
		"IP": "138.45.223.12", "txRef": "ACB-123", "payment_type": "account",
		"device_fingerprint": "675754758584e3847573",
	}

	response, _ = rave.ChargeAccount(accountDetails)

	v, _ := jason.NewObjectFromBytes(response)
	status, _ := v.GetString("status")

	assertEqual(t, status, "success")
}

// top up wallet
// The only way to top up Rave wallet while testing, is to do an Access Bank charge
// without funds in your account, the Refund test will fail
func topUpAccount() []byte {
	// get access bank details
	response, _ := rave.ListBanks()

	var banks []map[string]string
	json.Unmarshal(response, &banks)

	accessBankCode := banks[0]["bankcode"]

	accountDetails := map[string]interface{}{
		"accountnumber": "0690000031", "accountbank": accessBankCode, "currency": "NGN",
		"country": "NG", "amount": 5000, "email": "TestAccessBankCharge@gmail.com",
		"phonenumber": "08123456787", "firstname": "Access", "lastname": "AccessBank",
		"IP": "138.45.223.12", "txRef": "ACB-123", "payment_type": "account",
		"device_fingerprint": "675754758584e3847573",
	}

	response, _ = rave.ChargeAccount(accountDetails)

	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	transactionReference, _ := data.GetString("flwRef")

	// validate account charge
	transaction := map[string]interface{}{
		"transactionreference": transactionReference,
		"otp": "12345",
	}
	response, _ = rave.ValidateAccountCharge(transaction)

	return response
}

// Test Validate Account Charge for local banks (Access Bank)
func TestValidateAccountCharge(t *testing.T) {
	t.Parallel()

	response := topUpAccount()
	v, _ := jason.NewObjectFromBytes(response)
	successMessage, _ := v.GetString("message")
	data, _ := v.GetObject("data")
	chargedAmount, _ := data.GetInt64("charged_amount")

	if successMessage != "Charge Complete" || chargedAmount != 5000 {
		t.Error("Account Charge failed")
		fmt.Println(successMessage, chargedAmount)
	}

}

// Every method (That makes use of the MakePostRequest method)
// should return a response (as map[string]interface) for any failed request
// When the API returns an error
func TestErrorResponse(t *testing.T) {
	t.Parallel()

	// Make a request without including the cvv
	verveCard := map[string]interface{}{
		"name": "TestErrorResponse", "cardno": "5590131743294314", "currency": "NGN",
		"country": "NG", "amount": "300", "cvv": "887", "expirymonth": "11",
		"expiryyear": "20", "suggested_auth": "pin", "pin": "3310",
		"email": "TestErrorResponse@flutter.co", "firstname": "error", "lastname": "response",
		"phonenumber": "081245554343", "IP": "103.238.105.185", "txRef": "MXX-AYT-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
		"redirect_url":       "http://127.0.0.1",
	}
	_, err := rave.ChargeCard(verveCard)

	if err == nil {
		t.Error("'TestErrorResponse' didn't raise an error")
	}

	errorString := err.Error()
	if errorString != "Fraudulent. Transaction. Status Code: 400" {
		t.Errorf("Method didn't raise 'Fraudulent. Transaction. Status Code: 400' instead it raised %s", errorString)
	}
}

// We should get a list of all Nigerian banks we can charge
func TestListBanks(t *testing.T) {
	t.Parallel()

	response, _ := rave.ListBanks()

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
		"email":     "TestChargeCard@flutter.co",
		"firstname": "charge", "lastname": "card", "phonenumber": "081245554343",
		"IP": "103.238.105.185", "txRef": "MXX-AYT-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
		"redirect_url":       "http://127.0.0.1",
	}

	response, _ := rave.ChargeCard(masterCard)

	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	transactionReference, _ := data.GetString("flwRef")

	transaction := map[string]interface{}{
		"transaction_reference": transactionReference,
		"otp": "12345",
	}

	response, _ = rave.ValidateCharge(transaction)
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

// Use specific Preauth "test" keys
func setPreauthKeys() {
	// set Preauth keys for testing
	if rave.Live == false {
		os.Setenv("RAVE_PUBLICKEY", "FLWPUBK-8cd258c49f38e05292e5472b2b15906e-X")
		os.Setenv("RAVE_SECKEY", "FLWSECK-c51891678d48c39eff3701ff686bdb69-X")
	}
}

// Set Default auth keys
func setDefaultKeys() {
	// replace preauth keys with old keys
	os.Setenv("RAVE_PUBLICKEY", DefaultPublicKey)
	os.Setenv("RAVE_SECKEY", DefaultSecretKey)
}

// Test Preauth method
func TestPreauth(t *testing.T) {
	preauthMasterCard := map[string]interface{}{
		"name": "Preauth", "cardno": "5840406187553286", "currency": "NGN",
		"country": "NG", "cvv": "116", "amount": "300", "expiryyear": "18",
		"expirymonth": "09", "suggested_auth": "pin", "pin": "1111",
		"email": "preauth_void@flutter.co", "firstname": "preauth", "lastname": "preauth",
		"phonenumber": "081245554343", "IP": "103.238.105.185", "txRef": "MXX-AYT-4578",
		"device_fingerprint": "69e6b7f0sb722037ba8428b70fbe03986c",
		"redirect_url":       "http://127.0.0.1",
	}

	setPreauthKeys()
	defer setDefaultKeys()

	// Preauthorize card
	response, _ := rave.PreauthorizeCard(preauthMasterCard)
	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	status, _ := data.GetString("status")

	assertEqual(t, status, "pending-capture")
}

// Test Preauth => Capture
func TestPreauthCapture(t *testing.T) {
	preauthMasterCard := map[string]interface{}{
		"name": "PreauthCapture", "cardno": "5840406187553286", "currency": "NGN",
		"country": "NG", "cvv": "116", "amount": "300", "expiryyear": "18",
		"expirymonth": "09", "suggested_auth": "pin", "pin": "1111",
		"email": "preauth_void@flutter.co", "firstname": "preauth", "lastname": "capture",
		"phonenumber": "081245554343", "IP": "103.238.105.185", "txRef": "MXX-AYT-4578",
		"device_fingerprint": "69e6b7f0sb722037ba8428b70fbe03986c",
		"redirect_url":       "http://127.0.0.1",
	}

	// Pre authorize transaction
	setPreauthKeys()
	defer setDefaultKeys()

	response, _ := rave.PreauthorizeCard(preauthMasterCard)
	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	transactionReference, _ := data.GetString("flwRef")

	// Capture the amount
	response, _ = rave.Capture(map[string]interface{}{"flwRef": transactionReference})
	v, _ = jason.NewObjectFromBytes(response)
	successMessage, _ := v.GetString("message")

	assertEqual(t, successMessage, "Capture complete")

}

// Test Preauth => Capture => Refund
func TestPreauthCaptureRefund(t *testing.T) {
	preauthMasterCard := map[string]interface{}{
		"name": "PreauthCaptureRefund", "cardno": "5840406187553286", "currency": "NGN",
		"country": "NG", "cvv": "116", "amount": "300", "expiryyear": "18",
		"expirymonth": "09", "suggested_auth": "pin", "pin": "1111",
		"email": "preauth_void@flutter.co", "firstname": "preauth", "lastname": "capture_refund",
		"phonenumber": "081245554343", "IP": "103.238.105.185", "txRef": "MXX-AYT-4578",
		"device_fingerprint": "69e6b7f0sb722037ba8428b70fbe03986c",
		"redirect_url":       "http://127.0.0.1",
	}

	// Pre authorize transaction
	setPreauthKeys()
	defer setDefaultKeys()

	response, _ := rave.PreauthorizeCard(preauthMasterCard)
	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	transactionReference, _ := data.GetString("flwRef")

	// Capture the amount
	response, _ = rave.Capture(map[string]interface{}{"flwRef": transactionReference})
	v, _ = jason.NewObjectFromBytes(response)

	// Refund the amount
	transaction := map[string]interface{}{
		"ref": transactionReference, "action": "void",
	}
	response, _ = rave.RefundOrVoidPreauth(transaction)

	v, _ = jason.NewObjectFromBytes(response)
	successMessage, _ := v.GetString("message")

	assertEqual(t, successMessage, "Refund or void complete")
}

// Test Preauth => Void
func TestPreauthVoid(t *testing.T) {

	preauthMasterCard := map[string]interface{}{
		"name": "PreauthVoid", "cardno": "5840406187553286", "currency": "NGN",
		"country": "NG", "cvv": "116", "amount": "300", "expiryyear": "18",
		"expirymonth": "09", "suggested_auth": "pin", "pin": "1111",
		"email": "preauth_void@flutter.co", "firstname": "preauth", "lastname": "void",
		"phonenumber": "081245554343", "IP": "103.238.105.185", "txRef": "MXX-AYT-4578",
		"device_fingerprint": "69e6b7f0sb722037ba8428b70fbe03986c",
		"redirect_url":       "http://127.0.0.1",
	}

	// Preauthorize card
	setPreauthKeys()
	defer setDefaultKeys()

	response, _ := rave.PreauthorizeCard(preauthMasterCard)
	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	transactionReference, _ := data.GetString("flwRef")

	// Void the transaction
	transaction := map[string]interface{}{
		"ref": transactionReference, "action": "void",
	}
	response, _ = rave.RefundOrVoidPreauth(transaction)

	v, _ = jason.NewObjectFromBytes(response)
	successMessage, _ := v.GetString("message")

	assertEqual(t, successMessage, "Refund or void complete")
}

// Verify the status of a transaction
func TestVerifyTransaction(t *testing.T) {
	t.Parallel()

	// Initialize the transaction and get a valid transaction reference
	masterCard := map[string]interface{}{
		"name": "verifyTransaction", "cardno": "5438898014560229", "currency": "NGN",
		"country": "NG", "cvv": "789", "amount": "300", "expiryyear": "19",
		"expirymonth": "09", "suggested_auth": "pin", "pin": "3310",
		"email": "verifyTransaction@flutter.co", "firstname": "verify", "lastname": "transaction",
		"phonenumber": "081245554343", "IP": "103.238.105.185", "txRef": "MXX-AYT-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
		"redirect_url":       "http://127.0.0.1",
	}

	response, _ := rave.ChargeCard(masterCard)

	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	transactionReference, _ := data.GetString("flwRef")
	currency, _ := data.GetString("currency")

	// Pay for the transaction
	transaction := map[string]interface{}{
		"transaction_reference": transactionReference,
		"otp": "12345",
	}
	response, _ = rave.ValidateCharge(transaction)

	// Verify the transaction
	transaction = map[string]interface{}{
		"flw_ref": transactionReference, "normalize": "1",
		"currency": currency, "amount": "1000",
	}
	_, err := rave.VerifyTransaction(transaction)
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
		"email": "TestXrequery@flutter.co", "firstname": "xrequery", "lastname": "xrequery",
		"phonenumber": "081245554343", "IP": "103.238.105.185", "txRef": "MXX-AYT-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe031234e",
		"redirect_url":       "http://127.0.0.1",
	}

	response, _ := rave.ChargeCard(masterCard)

	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	transactionReference, _ := data.GetString("flwRef")
	currency, _ := data.GetString("currency")

	// Pay for the transaction
	transaction := map[string]interface{}{
		"transaction_reference": transactionReference,
		"otp": "12345",
	}
	response, _ = rave.ValidateCharge(transaction)

	// Verify the transaction
	// flw_ref is needed for verification
	transaction = map[string]interface{}{
		"flw_ref": transactionReference, "tx_ref": "abcdef",
		"last_attempt": "1", "only_attempt": "1",
		"currency": currency, "amount": "5300",
	}
	_, err := rave.XrequeryTransactionVerification(transaction)
	if err != nil {
		t.Error(err.Error())
	}
}

// Test transaction refund
func TestRefundTransaction(t *testing.T) {
	t.Parallel()

	// Initialize the transaction and get a valid transaction reference
	masterCard := map[string]interface{}{
		"name": "TestRefund", "cardno": "5438898014560229", "currency": "NGN",
		"country": "NG", "cvv": "789", "amount": "5300", "expiryyear": "19",
		"expirymonth": "09", "suggested_auth": "pin", "pin": "3310",
		"email": "TestRefundTransaction@flutter.co", "firstname": "refund",
		"lastname": "transaction", "phonenumber": "081245554343",
		"IP": "103.238.105.185", "txRef": "MXX-AYT-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbx031234e",
		"redirect_url":       "http://127.0.0.1",
	}

	response, _ := rave.ChargeCard(masterCard)

	v, _ := jason.NewObjectFromBytes(response)
	data, _ := v.GetObject("data")
	transactionReference, _ := data.GetString("flwRef")

	// Pay for the transaction
	transaction := map[string]interface{}{
		"transaction_reference": transactionReference,
		"otp": "12345",
	}
	response, _ = rave.ValidateCharge(transaction)

	// topup the Rave wallet
	topUpAccount()

	// Refund the transaction
	transaction = map[string]interface{}{"ref": transactionReference}
	response, _ = rave.RefundTransaction(transaction)
	v, _ = jason.NewObjectFromBytes(response)
	successMessage, _ := v.GetString("message")
	data, _ = v.GetObject("data")
	AmountRefunded, _ := data.GetInt64("AmountRefunded")

	if successMessage != "Refunded" || AmountRefunded != 5300 {
		t.Error("Transaction Refund failed")
		fmt.Println(successMessage, AmountRefunded)
	}
}

// Test Get fees endpoint
func testGetFees(t *testing.T) {
	t.Parallel()

	data := map[string]interface{}{
		"amount": "5300", "currency": "NGN",
	}

	response, _ := rave.GetFees(data)

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
	oldSecKey := rave.GetSecretKey()

	err := os.Setenv("RAVE_SECKEY", "FLWSECK-bb971402072265fb156e90a3578fe5e6-X")
	if err != nil {
		panic(err)
	}

	integrityChecksum := rave.CalculateIntegrityCheckSum(data)

	assertEqual(t, integrityChecksum, "a14ac4eba0902e8fd6b5fdf542f46d6efc18885a63c3d5f100c26715c7c8d8f4")

	// set "RAVE_SECKEY" to it's old value
	os.Setenv("RAVE_SECKEY", oldSecKey)
}

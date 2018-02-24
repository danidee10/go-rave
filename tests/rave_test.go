// Tests for the rave package

package tests

import (
	"fmt"
	"go-rave/rave"
	"reflect"
	"testing"
)

var Rave = rave.NewRave()

func assertEqual(t *testing.T, val1 interface{}, val2 interface{}) {
	if val1 != val2 {
		t.Fatalf(
			"'%s'(%s) is not Equal to '%s'(%s)",
			val1, reflect.TypeOf(val1), val2, reflect.TypeOf(val2),
		)
	}
}

// Setup test suite
func setUpTest(*testing.M) {
	Rave.Live = false
}

// Test the encryption function
func TestEncryption(t *testing.T) {
	seckey := Rave.GetSecretKey()
	encryptedSecretKey := Rave.GetKey(seckey)

	assertEqual(t, Rave.Encrypt3Des(encryptedSecretKey, "Hello world"), "fus4LnqrvKWXqm7wueoj2Q==")
}

// It should raise an error if the pin wasn't passed and the suggested_auth is "PIN"
func testSuggestedAuthRaisesError(t *testing.T) {
	masterCard := map[string]interface{}{
		"name": "hello", "cardno": "5438898014560229", "currency": "NGN",
		"country": "NG", "cvv": "789", "amount": "300", "expiryyear": "19",
		"expirymonth": "09",
		"email":       "tester@flutter.co", "IP": "103.238.105.185", "txRef": "MXX-ASC-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
	}
	ret := Rave.ChargeCard(masterCard)

	fmt.Println(ret)
}

func testSuggestedAuth(t *testing.T) {
	masterCard := map[string]interface{}{
		"name": "hello", "cardno": "5438898014560229", "currency": "NGN",
		"country": "NG", "cvv": "789", "amount": "300", "expiryyear": "19",
		"expirymonth": "09", "pin": "3310",
		"email": "tester@flutter.co", "IP": "103.238.105.185", "txRef": "MXX-ASC-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
	}
	ret := Rave.ChargeCard(masterCard)

	fmt.Println(ret)
}

func testMasterCardPaymentWithPin(t *testing.T) {
	masterCard := map[string]interface{}{
		"name": "hello", "cardno": "5438898014560229", "currency": "NGN",
		"country": "NG", "cvv": "789", "amount": "300", "expiryyear": "19",
		"expirymonth": "09", "suggested_auth": "pin", "pin": "3310",
		"email": "tester@flutter.co", "IP": "103.238.105.185", "txRef": "MXX-ASC-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
	}
	ret := Rave.ChargeCard(masterCard)

	fmt.Println(ret)
}

func testVerveCardPaymentWithPin(t *testing.T) {
	verveCard := map[string]interface{}{
		"name": "hello", "cardno": "5061020000000000094", "currency": "NGN",
		"country": "NG", "cvv": "347", "amount": "300", "expiryyear": "20",
		"expirymonth": "07", "suggested_auth": "pin", "pin": "1111",
		"email": "tester@flutter.co", "IP": "103.238.105.185", "txRef": "MXX-ASC-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
	}
	ret := Rave.ChargeCard(verveCard)

	fmt.Println(ret)
}

func testVisaPaymentWith3DSecure(t *testing.T) {
	visaCard := map[string]interface{}{
		"name": "hello", "cardno": "4187427415564246", "currency": "NGN",
		"country": "NG", "cvv": "828", "amount": "300", "expiryyear": "19",
		"expirymonth": "09", "email": "tester@flutter.co", "IP": "103.238.105.185",
		"txRef": "MXX-ASC-4578", "device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
	}
	ret := Rave.ChargeCard(visaCard)

	// authurl should be part of the response

	fmt.Println(ret)
}

// Every method (That makes use of the MakePostRequest method)
// should return a response (as map[string]interface) for any failed request
// When the API returns an error
func testErrorResponse(t *testing.T) {
	// Make a request without including the cvv
	verveCard := map[string]interface{}{
		"name": "hello", "cardno": "5061020000000000094", "currency": "NGN",
		"country": "NG", "amount": "300", "expiryyear": "20",
		"expirymonth": "07", "suggested_auth": "pin", "pin": "1111",
		"email": "tester@flutter.co", "IP": "103.238.105.185", "txRef": "MXX-ASC-4578",
		"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
	}
	ret := Rave.ChargeCard(verveCard)

	// should return a cvv required error

	fmt.Println(ret)
}

// We should get a list of all Nigerian banks we can charge
func testListBanks(t *testing.T) {
	response := Rave.ListBanks()

	fmt.Println(response)
}

// Test that a charge on a card can be validated using OTP
func testChargeCard(t *testing.T) {
	// TODO: make request and get a fresh transaction reference
	transaction := map[string]interface{}{
		"transaction_reference": "FLW-MOCK-539111aa99835cbbe028b058d4c9e961",
		"otp": "12345",
	}

	ret := Rave.ValidateCharge(transaction)

	fmt.Println(ret)
}

// Verify the status of a transaction
func testVerifyTransaction(t *testing.T) {
	// TODO: make sure we use a valid transaction reference
	transaction := map[string]interface{}{
		"flw_ref":   "FLW-MOCK-539111aa99835cbbe028b058d4c9e961",
		"normalize": "1",
	}

	ret := Rave.VerifyTransaction(transaction)

	fmt.Println(ret)
}

// Verify the status of a transaction using XRequery
func testXrequeryTransactionVerification(t *testing.T) {
	transaction := map[string]interface{}{
		"flw_ref": "FLW-MOCK-6f52518a2ecca2b6b090f9593eb390ce",
		"tx_ref":  "abcdef", "last_attempt": "1", "only_attempt": "1",
	}

	ret := Rave.XrequeryTransactionVerification(transaction)

	fmt.Println(ret)
}

// Test Get fees endpoint
func TestGetFees(t *testing.T) {
	data := map[string]interface{}{
		"amount": "1052.50", "currency": "NGN",
	}

	ret := Rave.GetFees(data)

	fmt.Println(ret)
}

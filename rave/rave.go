package rave

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/antonholmquist/jason"
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

// ChargeCard : Sends a Card request and determine the validation flow to be used
func (r rave) ChargeCard(chargeData map[string]interface{}) ([]byte, error) {
	err := checkRequiredParameters(chargeData, []string{"redirect_url"})
	if err != nil {
		return nil, err
	}

	postData := r.setUpCharge(chargeData)
	response, err := r.chargeCard(postData)
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

		postData = r.setUpCharge(chargeData)
		response, err = r.chargeCard(postData)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
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
func (r rave) chargeCard(data map[string]interface{}) ([]byte, error) {
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/charge"

	response, err := MakePostRequest(URL, data)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ValidateCharge : Validate a card charge using OTP
func (r rave) ValidateCharge(data map[string]interface{}) ([]byte, error) {

	data["PBFPubKey"] = r.GetPublicKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/validatecharge"

	response, err := MakePostRequest(URL, data)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ValidateAccountCharge : Validate an account charge using OTP
func (r rave) ValidateAccountCharge(data map[string]interface{}) ([]byte, error) {

	data["PBFPubKey"] = r.GetPublicKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/validate"

	response, err := MakePostRequest(URL, data)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// VerifyTransaction: Verify a transaction using "flw_ref" or "tx_ref"
func (r rave) VerifyTransaction(data map[string]interface{}) ([]byte, error) {
	err := checkRequiredParameters(data, []string{"amount", "currency", "flw_ref"})
	if err != nil {
		return nil, err
	}

	data["SECKEY"] = r.GetSecretKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/verify"

	response, err := MakePostRequest(URL, data)
	if err != nil {
		return nil, err
	}

	err = verifyTransaction(data, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r rave) XrequeryTransactionVerification(data map[string]interface{}) ([]byte, error) {
	err := checkRequiredParameters(data, []string{"amount", "currency", "flw_ref"})
	if err != nil {
		return nil, err
	}

	data["SECKEY"] = r.GetSecretKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/xrequery"

	response, err := MakePostRequest(URL, data)
	if err != nil {
		return nil, err
	}

	err = verifyTransaction(data, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Verify a transaction using the steps outlined in https://flutterwavedevelopers.readme.io/v1.0/reference#verification
func verifyTransaction(transactionData map[string]interface{}, response []byte) error {
	v, _ := jason.NewObjectFromBytes(response)
	successMessage, _ := v.GetString("message")

	data, _ := v.GetObject("data")
	transactionRef, err := data.GetString("flw_ref")
	if err != nil {
		// That means we're probably using XRequery search for "flwref" instead
		transactionRef, _ = data.GetString("flwref")
	}

	// declare differing variables
	var chargeResponse string
	var currency string
	var chargedAmount int64

	flwMeta, err := data.GetObject("flwMeta")
	if err != nil {
		// That means we're probably using XRequery get details from the "data" object instead
		chargeResponse, _ = data.GetString("chargecode")
		currency, _ = data.GetString("currency")
		chargedAmount, _ = data.GetInt64("chargedamount")

	} else {
		chargeResponse, _ = flwMeta.GetString("chargeResponse")
		currency, _ = data.GetString("transaction_currency")
		chargedAmount, _ = data.GetInt64("charged_amount")
	}

	// Run checks on the transaction
	transactionReference, _ := transactionData["flw_ref"]
	currencyCode := transactionData["currency"]

	if transactionRef != transactionReference {
		return fmt.Errorf("Transaction not verified because the transaction reference doesn't match: '%s' != '%s'", transactionRef, transactionReference)
	}

	if successMessage != "Tx Fetched" {
		return errors.New("Transaction not verified because success message is not equal to 'Tx Fetched'")
	}

	if chargeResponse != "00" && chargeResponse != "0" {
		return errors.New("Transaction not verified because the charged response is not equal to '00' or '0'")
	}

	if currency != currencyCode {
		return fmt.Errorf("Transaction not verified because the currency code doesn't match: '%s' != '%s'", currency, currencyCode)
	}

	amountString := transactionData["amount"].(string)
	amount, _ := strconv.ParseInt(amountString, 10, 32)
	if amount < chargedAmount {
		return errors.New("Transaction not verified, charged amount should be greater or equal amount to be paid")
	}

	return nil
}

// PreauthorizeCard : This is just a wrapper arond the Charge function that automatically
// sets "charge_type" to "pre_auth"
func (r rave) PreauthorizeCard(chargeData map[string]interface{}) ([]byte, error) {
	chargeData["charge_type"] = "pre_auth"

	response, err := r.ChargeCard(chargeData)
	if err != nil {
		return nil, err
	}

	return response, nil

}

func (r rave) Capture(data map[string]interface{}) ([]byte, error) {
	data["SECKEY"] = r.GetSecretKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/capture"

	response, err := MakePostRequest(URL, data)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// RefundOrVoid : Refund or void a captured amount
func (r rave) RefundOrVoid(data map[string]interface{}) ([]byte, error) {
	data["SECKEY"] = r.GetSecretKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/refundorvoid"

	response, err := MakePostRequest(URL, data)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetFees
func (r rave) GetFees(data map[string]interface{}) ([]byte, error) {
	data["PBFPubKey"] = r.GetPublicKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/fee"

	response, err := MakePostRequest(URL, data)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Refund : Refund direct charges
func (r rave) Refund(data map[string]interface{}) ([]byte, error) {
	data["seckey"] = r.GetSecretKey()
	URL := r.getBaseURL() + "/gpx/merchant/transactions/refund"

	response, err := MakePostRequest(URL, data)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ListBanks : List Nigerian banks.
func (r rave) ListBanks() ([]byte, error) {
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/flwpbf-banks.js?json=1"
	response, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	return body, nil
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

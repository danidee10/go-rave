/*
This file contains functions/methods that deals with transactions

Transactions can only be carried out on existing payments
*/

package rave

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/antonholmquist/jason"
)

// VerifyTransaction : Verify a transaction using "flw_ref" or "tx_ref"
func (r Rave) VerifyTransaction(data map[string]interface{}) ([]byte, error) {
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

// XrequeryTransactionVerification : verify a transaction using xrequery
func (r Rave) XrequeryTransactionVerification(data map[string]interface{}) ([]byte, error) {
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

	transactionReference, _ := transactionData["flw_ref"]
	currencyCode := transactionData["currency"]

	amountString := transactionData["amount"].(string)
	amount, _ := strconv.ParseInt(amountString, 10, 32)

	// Run "Five-Step" verification on the transaction
	err = verifyTransactionReference(transactionRef, transactionReference)
	err = verifySuccessMessage(successMessage)
	err = verifyChargeResponse(chargeResponse)
	err = verifyCurrencyCode(currency, currencyCode)
	err = verifyChargedAmount(chargedAmount, amount)
	if err != nil {
		return err
	}

	return nil
}

// The Transaction reference should match
func verifyTransactionReference(apiTransactionRef, funcTransactionRef interface{}) error {
	if apiTransactionRef != funcTransactionRef {
		return fmt.Errorf(
			"Transaction not verified because the transaction reference doesn't match: '%s' != '%s'",
			apiTransactionRef, funcTransactionRef,
		)
	}

	return nil
}

// The success message should equal "Tx Fetched" for a succesful transaction
func verifySuccessMessage(successMessage string) error {
	if successMessage != "Tx Fetched" {
		return errors.New("Transaction not verified because success message is not equal to 'Tx Fetched'")
	}

	return nil
}

// The Charge response should equal "00" or "0"
func verifyChargeResponse(chargeResponse string) error {
	if chargeResponse != "00" && chargeResponse != "0" {
		return errors.New("Transaction not verified because the charged response is not equal to '00' or '0'")
	}

	return nil
}

// The Currency code must match
func verifyCurrencyCode(apiCurrencyCode, funcCurrencyCode interface{}) error {
	if apiCurrencyCode != funcCurrencyCode {
		return fmt.Errorf(
			"Transaction not verified because the currency code doesn't match: '%s' != '%s'",
			apiCurrencyCode, funcCurrencyCode,
		)
	}

	return nil
}

// The Charged Amount must be greater than or equal to the paid amount
func verifyChargedAmount(apiChargedAmount, funcChargedAmount int64) error {
	if funcChargedAmount < apiChargedAmount {
		return errors.New("Transaction not verified, charged amount should be greater or equal amount to be paid")
	}

	return nil
}

// RefundTransaction : Refund direct charges
func (r Rave) RefundTransaction(data map[string]interface{}) ([]byte, error) {
	data["seckey"] = r.GetSecretKey()
	URL := r.getBaseURL() + "/gpx/merchant/transactions/refund"

	response, err := MakePostRequest(URL, data)
	if err != nil {
		return nil, err
	}

	return response, nil
}

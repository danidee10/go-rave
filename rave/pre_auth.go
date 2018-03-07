/* This file contains the functions/methods for Preauthorization */

package rave

// PreauthorizeCard : This is just a wrapper arond the ChargeCard method
// that automatically sets "charge_type" to "preauth"
func (r Rave) PreauthorizeCard(chargeData map[string]interface{}) ([]byte, error) {
	chargeData["charge_type"] = "preauth"

	response, err := r.ChargeCard(chargeData)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Capture : Capture a preauthorized transaction
func (r Rave) Capture(data map[string]interface{}) ([]byte, error) {
	err := checkRequiredParameters(data, []string{"flwRef"})
	if err != nil {
		return nil, err
	}

	data["SECKEY"] = r.GetSecretKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/capture"

	response, err := MakePostRequest(URL, data)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// RefundOrVoidPreauth : Refund or void a captured amount
func (r Rave) RefundOrVoidPreauth(data map[string]interface{}) ([]byte, error) {
	err := checkRequiredParameters(data, []string{"ref", "action"})

	data["SECKEY"] = r.GetSecretKey()
	URL := r.getBaseURL() + "/flwv3-pug/getpaidx/api/refundorvoid"

	response, err := MakePostRequest(URL, data)
	if err != nil {
		return nil, err
	}

	return response, nil
}

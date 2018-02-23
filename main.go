package main

import (
	"go-rave/rave"
)

func main() {
	r := rave.NewRave()
	r.Live = false

	/*
		response := r.ListBanks()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(body))

		==============================================
	*/

	/*
		seckey := r.GetSecretKey()
		encryptedSecretKey := r.GetKey(seckey)

		fmt.Println(r.Encrypt3Des(encryptedSecretKey, "hello world this is a long string"))

		==============================================
	*/

	/*

		masterCard := map[string]string{
			"name": "hello", "cardno": "5438898014560229", "currency": "NGN",
			"country": "NG", "cvv": "789", "amount": "300", "expiryyear": "19",
			"expirymonth": "09", "suggested_auth": "pin", "pin": "3310",
			"email": "tester@flutter.co", "IP": "103.238.105.185", "txRef": "MXX-ASC-4578",
			"device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
		}
		ret := r.ChargeCard(masterCard)

		fmt.Println(ret)

		==============================================
	*/

	/*
		transaction := map[string]string{
			"transaction_reference": "FLW-MOCK-539111aa99835cbbe028b058d4c9e961",
			"otp": "12345",
		}

		ret := r.ValidateCardCharge(transaction)
		fmt.Println(ret)
		===================================
	*/

	/*

		transaction := map[string]string{
			"flw_ref":   "FLW-MOCK-6f52518a2ecca2b6b090f9593eb390ce",
			"normalize": "1",
		}

		ret := r.VerifyTransaction(transaction)

		fmt.Println(ret)

		==============================================
	*/

	/*

		transaction := map[string]string{
			"flw_ref": "FLW-MOCK-6f52518a2ecca2b6b090f9593eb390ce",
			"tx_ref":  "abcdef", "last_attempt": "1", "only_attempt": "1",
		}

		ret := r.XrequeryTransactionVerification(transaction)

		fmt.Println(ret)

	*/

}

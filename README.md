# go-rave

[![Maintainability](https://api.codeclimate.com/v1/badges/60d8ae0dc97cbaca5089/maintainability)](https://codeclimate.com/github/danidee10/go-rave/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/60d8ae0dc97cbaca5089/test_coverage)](https://codeclimate.com/github/danidee10/go-rave/test_coverage)
<img class="badge" tag="github.com/danidee10/go-rave" src="/badge/github.com/danidee10/go-rave">

go-rave is a Go library for flutterwave's [rave](http://rave.frontendpwc.com/) payment platform.

It currently supports the following features:

* Account charge (NG Banks)

* Account charge (International for US and ZAR).

* Card Charge (Baked in support for 3DSecure/PIN).

* Encryption

* Transaction status check (Normal requery flow and xrequery).

* Retry transaction status check flow.

* Preauth -> Capture -> Refund/void flow.

* Support for USSD and Mcash (Alternative payment methods).

* List of banks for NG Account charge. (Get banks list).

* Get fees endpoint.

* Integrity Checksum (https://flutterwavedevelopers.readme.io/docs/checksum).

## Set Up

Go to [rave](http://ravepay.co/) and sign up.
This would provide you with a public and private authorization key which would be used throughout the library.

Store these authorization keys in your environment as `RAVE_PUBLICKEY` for the public key and `RAVE_SECKEY`.

They can be retrieved at runtime with `rave.GetPublicKey()` and `rave.GetSecretKey()` respectively.

## Getting Started

Install the library using `go get`

`go get github.com/danidee10/go-rave/rave`


After the installation, you can import the library like this:

``` go
import (
    ...
    "github.com/danidee10/go-rave/rave"
)
```

## Before you dive in!

go-rave was built with **version 2.0** of Rave's API in mind. If you need to check the documentation for any information, make sure you're looking at https://flutterwavedevelopers.readme.io/v2.0/reference#introduction

## Usage

All the functionality provided by the library is contained in the `Rave` `struct` to initialize a new instance of the `struct` call:

```go
Rave := rave.NewRave()
Rave.Live = false
```

Setting `Rave.Live` to `false` puts the library in development mode which means the base url is `http://flw-pms-dev.eu-west-1.elasticbeanstalk.com` when you're ready to go live, set `Rave.Live` to `true`

*Don't forget to update your publickey and seckey*

Since `Go` doesn't have keyword arguments most of the library's functions (that take input) use maps (`map[string]interface{}`).

For example, this is how you would represent a Master Card:

```go
masterCard := map[string]interface{}{
    "name": "hello", "cardno": "5438898014560229", "currency": "NGN",
    "country": "NG", "cvv": "789", "amount": "300", "expiryyear": "19",
    "expirymonth": "09", "suggested_auth": "pin", "pin": "3310",
    "email": "tester@flutter.co", "IP": "103.238.105.185", "txRef": "MXX-ASC-4578",
    "device_fingerprint": "69e6b7f0sb72037aa8428b70fbe03986c",
}
```

Parameters that are required by Rave's API are also checked. If any parameter is missing, `go-rave` will return an `error`.

All functions also return bytes. You can use a library like [jason](https://github.com/antonholmquist/jason) to read the values returned from the API.

## Library methods/functions

### Payment with card or account

**Documentation:** https://flutterwavedevelopers.readme.io/v1.0/reference#rave-direct-charge

```go
card := map[string]interface{}{..., "redirect_url": "http://127.0.0.1"}

response, err := Rave.ChargeCard(masterCard)
if err != nil {
    // handle error
    panic(err)
}
fmt.Println(response)
```

***Since it's not possible to determine the type of card (International or local) and the AuthModel required without consulting Rave's API, the 'redirect_url' parameter is mandatory for this function. You have to specify one so you can get the response back from Rave for an international card. This parameter isn't actually required for local cards.***

### Encrypting data

**Documentation:** https://flutterwavedevelopers.readme.io/v2.0/reference#rave-encryption

To encrypt data(card/account) with `3Des` call the `Rave.Encrypt3Des` method.

```go
encryptedData := Rave.Encrypt3Des(data)
```

***NOTE: You may not need to call this function if you use the methods provided by the library. The card/account data is automatically encrypted for you in any method that requires it.***


### Charge Validation

**Documentation:** https://flutterwavedevelopers.readme.io/v2.0/reference#rave-validation

Charge validation is handled by two methods.

#### Card

**Documentation:** https://flutterwavedevelopers.readme.io/v2.0/reference#validate-card-charge-1

**Required parameters:** `transaction_reference`, `otp`

To validate a charge for a card, call:

```go
transaction := map[string]interface{}{"transaction_reference": "...", "otp": "12345"}
response, err := Rave.ValidateCharge(transaction)
if err != nil {
    // handle error
}
```

with the transaction details.

#### Account

**Documentation:** https://flutterwavedevelopers.readme.io/v2.0/reference#validate-account-charge-1

**Required parameters:** `transactionreference`, `otp`

To validate a charge for an Account, call:

```go
transaction := map[string]interface{}{"transactionreference": "...", "otp": "12345"}
response, err := Rave.ValidateAccountCharge(transaction)
if err != nil {
    // handle error
}
```

with the transaction details.

### Transaction Verification

**Documentation:** https://flutterwavedevelopers.readme.io/v1.0/reference#verification

This is a very important function, because you have to make sure Rave has charged the user's card/account and received the payment before giving value to a user.

They're two ways of validating Rave transactions and `go-rave` allows you to use both. Each transaction is verified using the steps outlined in the [API documentation](https://flutterwavedevelopers.readme.io/v1.0/reference#verification). ***make sure you verify that no error was returned before giving value***

#### Normal Verification

**Documentation:** https://flutterwavedevelopers.readme.io/v2.0/reference#transaction-status-check

**Required parameters:** `flw_ref` (Transaction reference), `currency`, `amount`

To verify a transaction, call the `VerifyTransaction` method with the transaction details (reference etc) and handle any errors returned from the method.

```go
transaction := map[string]interface{}{
    ..., "flw_ref": transactionReference,
    "currency": currency, "amount": "1000",
}
response, err := Rave.VerifyTransaction(transaction)
if err != nil {
    // handle error || don't grant value
}
```

#### Transaction Verification with Xrequery

**Documentation:** https://flutterwavedevelopers.readme.io/v2.0/reference#xrequery-transaction-verification

**Required parameters:** `flw_ref` (Transaction reference), `currency`, `amount`

To verify a transaction with `Xrequery`, call the `XrequeryTransactionVerification` method with the transaction details.

```go
transaction := map[string]interface{}{
    ..., "flw_ref": transactionReference,
    "currency": currency, "amount": "1000",
}
response, err := Rave.XrequeryTransactionVerification(transaction)
if err != nil {
    // handle error || don't grant value
}
```

#### Refund

**Documentation:** https://flutterwavedevelopers.readme.io/v2.0/reference#refund

**Required parameters:** `ref` (Transaction reference)

To Initiate a refund call the `Refund` method with the transaction details.

```go
transaction := map[string]interface{}{"ref": "..."}
response, err := Rave.Refund(transaction)
if err != nil {
    // handle error
}
```

### List of Banks

**Documentation:** https://flutterwavedevelopers.readme.io/v2.0/reference#list-of-banks

Simply call `ListBanks`:

```go
response, err := Rave.ListBanks()
if err != nil {
    // handle error
}
```

### Get Fees

**Documentation:** https://flutterwavedevelopers.readme.io/v2.0/reference#get-fees

**Required parameters:** `amount`, `currency`

To get the fee for a particular amount call `GetFee` with valid details.

```go
data := map[string]interface{}{
    "amount": "1052.50", "currency": "NGN",
}
response, err := Rave.GetFees(data)
if err != nil {
    // handle error
}
```

### Preauthorize card

**Documentation:** https://flutterwavedevelopers.readme.io/v2.0/reference#preauthorise-card

To Preauthorize a card, call the `PreAuthorizeCard` method.

***NOTE: The client data is encrypted automatically and `charge_type` is also set to `preauth` for you.***

This is a simple code snippet showing you how to achieve that

```go
card := map[string]interface{}{...}
response, err := Rave.PreauthorizeCard(card)
if err != nil {
    // handle error
}
```

### Preauthorization Capture

**Documentation:** https://flutterwavedevelopers.readme.io/v2.0/reference#capture

**Required parameters:** `flwRef` (Transaction refernece).

To Capture an amount, call the `Capture` method with valid data (Typically the response gotten from `Rave.PreauthorizeCard`)

```go
transaction := map[string]interface{}{"flwRef": "..."}
response, err := Rave.Capture(transaction)
if err != nil {
    // handle error
}
```

### Transaction Refund or Void

**Documentation:** https://flutterwavedevelopers.readme.io/v2.0/reference#refund-or-void

Finally, To Refund or void a transaction, call the `RefundOrVoid` method with valid transaction data (containing the `reference`)

```go
transaction := map[string]interface{}{...}
response, err := Rave.RefundOrVoid(transaction)
if err != nil {
    // handle error
}
```

### IntegrityCheckSum

The Integrity checksum is necessary to secure payments on the client side. To generate an integrity hash call the `CalculateIntegrityCheckSum` and pass in the data.

```go
data := map[string]interface{}{...}
integrityCheckSum := Rave.CalculateIntegrityCheckSum(data)
```

## Contributing

To contribute, fork the repo, make your changes, write tests (If necessary) and create a pull request.

## Todo

 Add More Tests

## Authors

* [Osaetin Daniel](https://github.com/danidee10)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

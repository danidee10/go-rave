# go-rave

go-rave is a python library for the flutterwave's [rave](http://rave.frontendpwc.com/) payment platform.

It currently supports the following features:

* Account charge (NG Banks)

* Account charge (International for US and ZAR).

* Card Charge (Bake in support for 3DSecure/PIN).

* Encryption

* Transaction status check (Normal requery flow and xrequery).

* Retry transaction status check flow.

* Preauth -> Capture -> Refund/void.

* Support for USSD and Mcash (Alternative payment methods).

* List of banks for NG Account charge. (Get banks list).

* Get fees endpoint.

* Integrity Checksum (https://flutterwavedevelopers.readme.io/docs/checksum).

## Set Up

Go to [rave](http://rave.frontendpwc.com/) and sign up.
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

Parameters that are required by Rave's API are also checked. If any parameter is missing, `go-rave` will spit out a helpful log message and crash.

All functions also return maps. To make it easier to read responses from the API.

### Payment with card or account

```go
card := map[string]interface{}{...}

response := Rave.ChargeCard(masterCard)
fmt.Println(response)
```

### Encrypting data

To encrypt data(card/account) with `3Des` call the `Rave.3D`

```go
encryptedData := Rave.Encrypt3Des(data)
```

***NOTE: You may not need to call this function if you use the methods provided by the library. The card/account data is automatically encrypted for you in any method that requires it.***


### Charge Validation

Charge validation is handled by two methods.

#### Card

To validate a charge for a card, call:

```go
transaction := map[string]interface{}{...}
response:= Rave.ValidateCharge(transaction)
```

with the transaction details.

#### Account

To validate a charge for an account, call:

```go
transaction := map[string]interface{}{...}
response := Rave.ValidateAccountCharge(transaction)
```

with the transaction details.

### Transaction Verification

To verify a transaction, call the `VerifyTransaction` method with the transaction details (reference etc).

```go
transaction := map[string]interface{}{...}
response := Rave.VerifyTransaction(transaction)
```

### Transaction Verification with Xrequery

To verify a transaction with `Xrequery`, call the `XrequeryTransactionVerification` method with the transaction details.

```go
transaction := map[string]interface{}{...}
Rave.XrequeryTransactionVerification(transaction)
```

#### Refund

To Initiate a refund call the `Refund` method with the transaction details.

```go
transaction := map[string]interface{}{...}
response := Rave.Refund(transaction)
```

### List of Banks

Simply call `ListBanks`:

```go
response := Rave.ListBanks()
```

### List Fees

To get the fee for a particular amount call `GetFee` with valid details.

```go
data := map[string]interface{}{
    "amount": "1052.50", "currency": "NGN",
}
response := Rave.GetFees(data)
```

### Preauthorize card

To Preauthorize a card, call the `PreAuthorizeCard` method.

***NOTE: The client data is encrypted for you automatically and `charge_type` is also set to `pre_auth` for you.***

This is a simple code snippet showing you how to achieve that

```go
card := map[string]interface{}{...}
response := Rave.PreauthorizeCard(card)
```

### Preauthorization Capture

To Capture an amount, call the `Capture` method with valid data (Typically the response gotten from `Rave.PreauthorizeCard`)

```go
transaction := map[string]interface{}{...}
response := Rave.Capture(transaction)
```

### Transaction Refund or Void

Finally, To Refund or void a transaction, call the `RefundOrVoid` method with valid transaction data (containing the `reference`)

```go
transaction := map[string]interface{}{...}
response := Rave.RefundOrVoid(transaction)
```

## Contributing

To contribute, fork the repo, make your changes, write tests (If necessary) and create a pull request.

## Todo

 Add More Tests

## Authors

* [Osaetin Daniel](https://github.com/danidee10)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

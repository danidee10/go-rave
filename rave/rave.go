package rave

import (
	"log"
	"os"
)

// Rave : Base Rave type
type Rave struct {
	Live    bool
	liveURL string
	testURL string

	publicKey string
	secretKey string
}

// getBaseURL : Returns the Correct URL based on Live status.
func (r Rave) getBaseURL() string {
	if r.Live {
		return r.liveURL
	}

	return r.testURL
}

// GetPublicKey : Get Rave Public key
func (r Rave) GetPublicKey() string {
	publicKey, found := os.LookupEnv("RAVE_PUBLICKEY")
	if !found {
		log.Fatal("You must set the \"RAVE_PUBLICKEY\" environment variable")
	}

	return publicKey
}

// GetSecretKey : Get Rave Secret key
func (r Rave) GetSecretKey() string {
	secKey, found := os.LookupEnv("RAVE_SECKEY")
	if !found {
		log.Fatal("You must set the \"RAVE_SECKEY\" environment variable")
	}

	return secKey
}

// NewRave : Constructor for Rave struct
func NewRave() Rave {
	Rave := Rave{}
	Rave.testURL = "http://flw-pms-dev.eu-west-1.elasticbeanstalk.com"
	Rave.liveURL = "https://api.ravepay.co"

	// default mode is development
	Rave.Live = false

	return Rave
}

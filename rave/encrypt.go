// Implements Rave Encryption Algorithm

package rave

import (
	"bytes"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

// getKey : Get a key for encryption
func (r Rave) getKey(seckey string) string {
	hashedSeckey := md5.Sum([]byte(seckey))
	hashedSeckeyLast12 := hashedSeckey[len(hashedSeckey)-6:] // -6 because it's a hex byte array not a string
	seckeyAdjusted := strings.Replace(seckey, "FLWSECK-", "", 1)
	seckeyAdjustedFirst12 := seckeyAdjusted[:12]

	return seckeyAdjustedFirst12 + hex.EncodeToString(hashedSeckeyLast12[:])
}

// pkcs5Padding : Implements PKCS5 padding
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// Encrypt3Des : Encrypts the data using 3Des encryption
// Go doesn't include ECB encryption in the standard library for security reasons
// reference: https://gist.github.com/cuixin/10612934
func (r Rave) Encrypt3Des(payload string) string {
	seckey := r.GetSecretKey()
	encryptedSecretKey := r.getKey(seckey)

	return r.encrypt3Des(encryptedSecretKey, payload)
}

// contains the logic for encryption with 3Des
func (r Rave) encrypt3Des(key string, payload string) string {
	block, err := des.NewTripleDESCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	bs := block.BlockSize() // block size is 8 by default
	payloadBytes := pkcs5Padding([]byte(payload), bs)

	if len(payloadBytes)%bs != 0 {
		panic("Need a multiple of the blocksize")
	}
	encrypted := make([]byte, len(payloadBytes))
	dst := encrypted

	for len(payloadBytes) > 0 {
		block.Encrypt(dst, payloadBytes[:bs])
		payloadBytes = payloadBytes[bs:]
		dst = dst[bs:]
	}

	return base64.StdEncoding.EncodeToString(encrypted)
}

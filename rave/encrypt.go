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

// GetKey : Get a key for encryption
func (r rave) GetKey(seckey string) string {
	hashedSeckey := md5.Sum([]byte(seckey))
	hashedSeckeyLast12 := hashedSeckey[len(hashedSeckey)-6:] // -6 because it's a hex byte array not a string
	seckeyAdjusted := strings.Replace(seckey, "FLWSECK-", "", 1)
	seckeyAdjustedFirst12 := seckeyAdjusted[:12]

	return seckeyAdjustedFirst12 + hex.EncodeToString(hashedSeckeyLast12[:])
}

// PKCS5Padding : Implements PKCS5 padding
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// Encrypt3Des : Encrypts the data using 3Des encryption
// Go doesn't include ECB encryption in the standard library for security reasons
// reference: https://gist.github.com/cuixin/10612934
func (r rave) Encrypt3Des(key string, payload string) string {
	block, err := des.NewTripleDESCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	bs := block.BlockSize() // block size is 8 by default
	payloadBytes := PKCS5Padding([]byte(payload), bs)

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

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"log"
)

// function to encrypt message to be sent
func encrypt(msg string, key rsa.PublicKey) string {

	label := []byte("OAEP Encrypted")
	rng := rand.Reader

	// * using OAEP algorithm to make it more secure
	// * using sha256
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &key, []byte(msg), label)
	// check for errors
	if err != nil {
		log.Fatalln("unable to encrypt")
	}

	return base64.StdEncoding.EncodeToString(ciphertext)
}

// function to decrypt message to be received
func decrypt(cipherText string, key rsa.PrivateKey) string {

	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")
	rng := rand.Reader

	// decrypting based on same parameters as encryption
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &key, ct, label)
	// check for errors
	if err != nil {
		log.Fatalln(err)
	}
	return string(plaintext)
}

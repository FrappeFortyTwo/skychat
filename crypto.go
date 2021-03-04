package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"log"
)

func encrypt(msg string, key rsa.PublicKey) string {

	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &key, []byte(msg), label)

	if err != nil {
		log.Fatalln("unable to encrypt")
	}

	return base64.StdEncoding.EncodeToString(ciphertext)
}

func decrypt(cipherText string, key rsa.PrivateKey) string {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &key, ct, label)
	if err != nil {
		log.Fatalln(err)
	}
	return string(plaintext)
}

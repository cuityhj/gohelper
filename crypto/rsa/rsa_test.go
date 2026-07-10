package rsa

import (
	"bytes"
	"os"
	"testing"

	ut "github.com/cuityhj/cement/unittest"
)

func TestRSA(t *testing.T) {
	privKey, err := GenerateRSAKeyPair(2048)
	ut.Assert(t, err == nil, "generate private key failed")

	pubKey := &privKey.PublicKey

	password := []byte("admin@123")
	encryptedPassword, err := Encrypt(pubKey, password)
	ut.Assert(t, err == nil, "encrypt password with public key failed")
	decryptedPassword, err := Decrypt(privKey, encryptedPassword)
	ut.Assert(t, err == nil, "decrypt password with private key failed")
	ut.Assert(t, bytes.Equal(decryptedPassword, password), string(password))

	privKeyFileName := "private.pem"
	err = WritePrivateKeyToFile(privKey, privKeyFileName)
	ut.Assert(t, err == nil, "write private key failed")
	defer os.Remove(privKeyFileName)

	pubKeyFileName := "public.pem"
	err = WritePublicKeyToFile(pubKey, pubKeyFileName)
	ut.Assert(t, err == nil, "write public key failed")
	defer os.Remove(pubKeyFileName)

	privKey, err = LoadPrivateKeyFromFile(privKeyFileName)
	ut.Assert(t, err == nil, "load private key failed")

	pubKey, err = LoadPublicKeyFromFile(pubKeyFileName)
	ut.Assert(t, err == nil, "load public key failed")

	encryptedPassword, err = Encrypt(pubKey, password)
	ut.Assert(t, err == nil, "encrypt password with public key failed")
	decryptedPassword, err = Decrypt(privKey, encryptedPassword)
	ut.Assert(t, err == nil, "decrypt password with private key failed")
	ut.Assert(t, bytes.Equal(decryptedPassword, password), string(password))
}

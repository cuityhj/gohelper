package ecies

import (
	"bytes"
	"os"
	"testing"

	ut "github.com/cuityhj/cement/unittest"
)

func TestECIES(t *testing.T) {
	privKey, pubKey, err := GenerateECIESKeyPair()
	if err != nil {
		ut.Assert(t, err == nil, "generate ecies key failed")
	}

	plaintext := []byte("this is a test plaintext for verifying ECIES encryption and decryption")
	ciphertext, err := pubKey.Encrypt(plaintext)
	ut.Assert(t, err == nil, "ecies encrypt failed")
	decrypted, err := privKey.Decrypt(ciphertext)
	ut.Assert(t, err == nil, "ecies decrypt failed")
	ut.Assert(t, bytes.Equal(decrypted, plaintext), "")

	privKeyFileName := "private.pem"
	err = privKey.WriteToFile(privKeyFileName)
	ut.Assert(t, err == nil, "write private key failed")
	defer os.Remove(privKeyFileName)

	pubKeyFileName := "public.pem"
	err = pubKey.WriteToFile(pubKeyFileName)
	ut.Assert(t, err == nil, "write public key failed")
	defer os.Remove(pubKeyFileName)

	privKey, err = LoadPrivateKeyFromFile(privKeyFileName)
	ut.Assert(t, err == nil, "load private key failed")

	pubKey, err = LoadPublicKeyFromFile(pubKeyFileName)
	ut.Assert(t, err == nil, "load public key failed")

	ciphertext, err = pubKey.Encrypt(plaintext)
	ut.Assert(t, err == nil, "encrypt password with public key failed")
	decrypted, err = privKey.Decrypt(ciphertext)
	ut.Assert(t, err == nil, "decrypt password with private key failed")
	ut.Assert(t, bytes.Equal(decrypted, plaintext), "")
}

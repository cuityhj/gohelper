package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func ReadPublicKeyFromFile(fileName string) (*rsa.PublicKey, error) {
	pubKeyPEM, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read public key failed: %s", err.Error())
	}

	return ParsePublicKey(pubKeyPEM)
}

func ParseBase64PublicKey(base64PubKey string) (*rsa.PublicKey, error) {
	if pubKeyPEM, err := base64.StdEncoding.DecodeString(base64PubKey); err != nil {
		return nil, fmt.Errorf("base64 decode public key failed: %s", err.Error())
	} else {
		return ParsePublicKey(pubKeyPEM)
	}
}

func ParsePublicKey(pubKeyPEM []byte) (*rsa.PublicKey, error) {
	publicKey, err := parsePublicKey(pubKeyPEM)
	if err != nil {
		return nil, err
	}

	if publicKey.N.BitLen() < 2048 {
		return nil, fmt.Errorf("rsa public key is smaller than 2048 bits")
	}

	return publicKey, nil
}

func parsePublicKey(pubKeyPEM []byte) (*rsa.PublicKey, error) {
	pubBlock, _ := pem.Decode(pubKeyPEM)
	if pubBlock == nil {
		return nil, fmt.Errorf("can not decode public key PEM")
	}

	switch pubBlock.Type {
	case "RSA PUBLIC KEY":
		return x509.ParsePKCS1PublicKey(pubBlock.Bytes)
	case "PUBLIC KEY":
		pubKey, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("x509 parse public key failed: %s", err.Error())
		}

		publicKey, ok := pubKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("public key is not rsa key")
		}

		return publicKey, nil
	default:
		return nil, fmt.Errorf("unsupported public key type %s", pubBlock.Type)
	}
}

func EncryptWithPublicPEM(plaintext string, pubKeyPEM []byte) (string, error) {
	publicKey, err := ParsePublicKey(pubKeyPEM)
	if err != nil {
		return "", err
	}

	return EncryptWithPublicKey(plaintext, publicKey)
}

func EncryptWithPublicKey(plaintext string, publicKey *rsa.PublicKey) (string, error) {
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(plaintext))
	if err != nil {
		return "", fmt.Errorf("rsa public key encode failed: %s", err.Error())
	}

	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

func ReadPrivateKeyFromFile(fileName string) (*rsa.PrivateKey, error) {
	privKeyPEM, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read private key failed: %s", err.Error())
	}

	return ParsePrivateKey(privKeyPEM)
}

func ParseBase64PrivateKey(bas64PrivKey string) (*rsa.PrivateKey, error) {
	if privateKeyPEM, err := base64.StdEncoding.DecodeString(bas64PrivKey); err != nil {
		return nil, fmt.Errorf("base64 decode private key failed: %s", err.Error())
	} else {
		return ParsePrivateKey(privateKeyPEM)
	}
}

func ParsePrivateKey(privKeyPEM []byte) (*rsa.PrivateKey, error) {
	privateKey, err := parsePrivateKey(privKeyPEM)
	if err != nil {
		return nil, err
	}

	if privateKey.N.BitLen() < 2048 {
		return nil, fmt.Errorf("rsa private key is smaller than 2048 bits")
	}

	return privateKey, nil
}

func parsePrivateKey(privKeyPEM []byte) (*rsa.PrivateKey, error) {
	privBlock, _ := pem.Decode(privKeyPEM)
	if privBlock == nil {
		return nil, fmt.Errorf("can not decode private key PEM")
	}

	switch privBlock.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	case "PRIVATE KEY":
		privKey, err := x509.ParsePKCS8PrivateKey(privBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("x509 parse private key failed: %s", err.Error())
		}

		privateKey, ok := privKey.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is not rsa key")
		}

		return privateKey, nil
	default:
		return nil, fmt.Errorf("unsupported private key type %s", privBlock.Type)
	}
}

func DecryptWithPrivatePEM(encryptedData string, privKeyPEM []byte) ([]byte, error) {
	privateKey, err := ParsePrivateKey(privKeyPEM)
	if err != nil {
		return nil, err
	}

	return DecryptWithPrivateKey(encryptedData, privateKey)
}

func DecryptWithPrivateKey(encryptedData string, privateKey *rsa.PrivateKey) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("base64 decode encrypted data failed: %s", err.Error())
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("rsa private key decrypt data failed: %s", err.Error())
	}

	return plaintext, nil
}

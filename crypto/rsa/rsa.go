package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

const (
	PKCS1_PUBLIC_KEY_TYPE  = "RSA PUBLIC KEY"
	PKIX_PUBLIC_KEY_TYPE   = "PUBLIC KEY"
	PKCS1_PRIVATE_KEY_TYPE = "RSA PRIVATE KEY"
	PKCS8_PRIVATE_KEY_TYPE = "PRIVATE KEY"
)

type (
	PrivateKey rsa.PrivateKey
	PublicKey  rsa.PublicKey
)

func GenerateRSAKeyPair(keyBits int) (*PrivateKey, *PublicKey, error) {
	if privKey, err := rsa.GenerateKey(rand.Reader, keyBits); err != nil {
		return nil, nil, err
	} else {
		return (*PrivateKey)(privKey), (*PublicKey)(&privKey.PublicKey), nil
	}
}

func (pub *PublicKey) WriteToFile(fileName string) error {
	pubKeyFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("create public key file failed: %s", err.Error())
	}

	defer pubKeyFile.Close()

	pubKeyDER, err := x509.MarshalPKIXPublicKey((*rsa.PublicKey)(pub))
	if err != nil {
		return fmt.Errorf("x509 marshal public key failed: %s", err.Error())
	}

	pubKeyPEM := &pem.Block{
		Type:  PKIX_PUBLIC_KEY_TYPE,
		Bytes: pubKeyDER,
	}
	if err := pem.Encode(pubKeyFile, pubKeyPEM); err != nil {
		return fmt.Errorf("write public key to file failed: %s", err.Error())
	}

	return nil
}

func LoadPublicKeyFromFile(fileName string) (*PublicKey, error) {
	pubKeyPEM, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read rsa public key failed: %s", err.Error())
	}

	return ParsePublicKey(pubKeyPEM)
}

func ParseBase64PublicKey(base64PubKey string) (*PublicKey, error) {
	if pubKeyPEM, err := base64.StdEncoding.DecodeString(base64PubKey); err != nil {
		return nil, fmt.Errorf("base64 decode public key failed: %s", err.Error())
	} else {
		return ParsePublicKey(pubKeyPEM)
	}
}

func ParsePublicKey(pubKeyPEM []byte) (*PublicKey, error) {
	if len(pubKeyPEM) == 0 {
		return nil, fmt.Errorf("public pem cannot be empty")
	}

	pubKeyBlock, _ := pem.Decode(pubKeyPEM)
	if pubKeyBlock == nil {
		return nil, fmt.Errorf("cannot decode public key PEM")
	}

	switch pubKeyBlock.Type {
	case PKCS1_PUBLIC_KEY_TYPE:
		pubKey, err := x509.ParsePKCS1PublicKey(pubKeyBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("x509 parse public key failed: %s", err.Error())
		}

		return (*PublicKey)(pubKey), nil
	case PKIX_PUBLIC_KEY_TYPE:
		publicKey, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("x509 parse public key failed: %s", err.Error())
		}

		pubKey, ok := publicKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("public key is not rsa key")
		}

		return (*PublicKey)(pubKey), nil
	default:
		return nil, fmt.Errorf("unsupported public key type %s", pubKeyBlock.Type)
	}
}

func EncryptWithPublicPEM(plaintext string, pubKeyPEM []byte) (string, error) {
	if len(plaintext) == 0 || len(pubKeyPEM) == 0 {
		return "", fmt.Errorf("plaintext or public pem cannot be empty")
	}

	pubKey, err := ParsePublicKey(pubKeyPEM)
	if err != nil {
		return "", err
	}

	return pubKey.Encrypt([]byte(plaintext))
}

func (pub *PublicKey) Encrypt(plaintext []byte) (string, error) {
	if len(plaintext) == 0 {
		return "", fmt.Errorf("plaintext cannot be empty")
	}

	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, (*rsa.PublicKey)(pub), plaintext)
	if err != nil {
		return "", fmt.Errorf("rsa public key encode failed: %s", err.Error())
	}

	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

func (priv *PrivateKey) Public() *PublicKey {
	return (*PublicKey)(&priv.PublicKey)
}

func (priv *PrivateKey) WriteToFile(fileName string) error {
	privKeyFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("create private key file failed: %s", err.Error())
	}

	defer privKeyFile.Close()

	privKeyDER, err := x509.MarshalPKCS8PrivateKey((*rsa.PrivateKey)(priv))
	if err != nil {
		return fmt.Errorf("x509 marshal private key failed: %s", err.Error())
	}

	privKeyPEM := &pem.Block{
		Type:  PKCS8_PRIVATE_KEY_TYPE,
		Bytes: privKeyDER,
	}
	if err := pem.Encode(privKeyFile, privKeyPEM); err != nil {
		return fmt.Errorf("write private key to file failed: %s", err.Error())
	}

	return nil
}

func LoadPrivateKeyFromFile(fileName string) (*PrivateKey, error) {
	privKeyPEM, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read private key failed: %s", err.Error())
	}

	return ParsePrivateKey(privKeyPEM)
}

func ParseBase64PrivateKey(bas64PrivKey string) (*PrivateKey, error) {
	if privKeyPEM, err := base64.StdEncoding.DecodeString(bas64PrivKey); err != nil {
		return nil, fmt.Errorf("base64 decode private key failed: %s", err.Error())
	} else {
		return ParsePrivateKey(privKeyPEM)
	}
}

func ParsePrivateKey(privKeyPEM []byte) (*PrivateKey, error) {
	if len(privKeyPEM) == 0 {
		return nil, fmt.Errorf("private pem cannot be empty")
	}

	privKeyBlock, _ := pem.Decode(privKeyPEM)
	if privKeyBlock == nil {
		return nil, fmt.Errorf("cannot decode private key PEM")
	}

	switch privKeyBlock.Type {
	case PKCS1_PRIVATE_KEY_TYPE:
		privKey, err := x509.ParsePKCS1PrivateKey(privKeyBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("x509 parse private key failed: %s", err.Error())
		}

		return (*PrivateKey)(privKey), nil
	case PKCS8_PRIVATE_KEY_TYPE:
		privateKey, err := x509.ParsePKCS8PrivateKey(privKeyBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("x509 parse private key failed: %s", err.Error())
		}

		privKey, ok := privateKey.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is not rsa key")
		}

		return (*PrivateKey)(privKey), nil
	default:
		return nil, fmt.Errorf("unsupported private key type %s", privKeyBlock.Type)
	}
}

func DecryptWithPrivatePEM(encryptedData string, privKeyPEM []byte) ([]byte, error) {
	if len(encryptedData) == 0 || len(privKeyPEM) == 0 {
		return nil, fmt.Errorf("ciphertext or private pem cannot be empty")
	}

	privKey, err := ParsePrivateKey(privKeyPEM)
	if err != nil {
		return nil, err
	}

	return privKey.Decrypt(encryptedData)
}

func (priv *PrivateKey) Decrypt(encryptedData string) ([]byte, error) {
	if len(encryptedData) == 0 {
		return nil, fmt.Errorf("ciphertext cannot be empty")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("base64 decode encrypted data failed: %s", err.Error())
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, (*rsa.PrivateKey)(priv), ciphertext)
	if err != nil {
		return nil, fmt.Errorf("rsa private key decrypt data failed: %s", err.Error())
	}

	return plaintext, nil
}

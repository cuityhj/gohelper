package ecies

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"

	"golang.org/x/crypto/hkdf"
)

const (
	PKCS8_PRIVATE_KEY_TYPE = "PRIVATE KEY"
	PKIX_PUBLIC_KEY_TYPE   = "PUBLIC KEY"

	AES_KEY_LENGTH   = 16
	GCM_NONCE_LENGTH = 12

	COMPRESSED_PUB_LENGTH = 33
)

var HKDF_INFO = []byte("ecies-p256-aes128gcm")

func GenerateECIESKeyPair() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func WritePrivateKeyToFile(privKey *ecdsa.PrivateKey, fileName string) error {
	privKeyFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("create private key file failed: %s", err.Error())
	}

	defer privKeyFile.Close()

	privKeyDER, err := x509.MarshalPKCS8PrivateKey(privKey)
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

func WritePublicKeyToFile(pubKey *ecdsa.PublicKey, fileName string) error {
	pubKeyFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("create public key file failed: %s", err.Error())
	}

	defer pubKeyFile.Close()

	pubKeyDER, err := x509.MarshalPKIXPublicKey(pubKey)
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

func LoadPrivateKeyFromFile(fileName string) (*ecdsa.PrivateKey, error) {
	privKeyPEM, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read private key failed: %s", err.Error())
	}

	return ParsePrivateKey(privKeyPEM)
}

func ParsePrivateKey(privKeyPEM []byte) (*ecdsa.PrivateKey, error) {
	privKeyBlock, _ := pem.Decode(privKeyPEM)
	if privKeyBlock == nil || privKeyBlock.Type != PKCS8_PRIVATE_KEY_TYPE {
		return nil, fmt.Errorf("can not decode private key PEM")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(privKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509 parse private key failed: %s", err.Error())
	}

	privKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not ecdsa key")
	}

	return privKey, nil
}

func LoadPublicKeyFromFile(fileName string) (*ecdsa.PublicKey, error) {
	pubKeyPEM, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read public key failed: %s", err.Error())
	}

	return ParsePublicKey(pubKeyPEM)
}

func ParsePublicKey(pubKeyPEM []byte) (*ecdsa.PublicKey, error) {
	pubKeyBlock, _ := pem.Decode(pubKeyPEM)
	if pubKeyBlock == nil || pubKeyBlock.Type != PKIX_PUBLIC_KEY_TYPE {
		return nil, fmt.Errorf("can not decode public key PEM")
	}

	publicKey, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509 parse public key failed: %s", err.Error())
	}

	pubKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not ecdsa key")
	}

	return pubKey, nil
}

func Encrypt(pubKey *ecdsa.PublicKey, plaintext []byte) (string, error) {
	curve := elliptic.P256()
	ephemeralPriv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return "", fmt.Errorf("ecdsa generate private key failed: %s", err.Error())
	}

	sharedSecret, err := ecdhDeriveSharedSecret(ephemeralPriv, pubKey)
	if err != nil {
		return "", err
	}

	gcm, nonce, err := newAESGCMAndNONCE(sharedSecret)
	if err != nil {
		return "", err
	}

	symCiphertext := gcm.Seal(nil, nonce, plaintext, nil)
	compressedPub := elliptic.MarshalCompressed(curve, ephemeralPriv.PublicKey.X, ephemeralPriv.PublicKey.Y)
	result := make([]byte, len(compressedPub)+len(symCiphertext))
	copy(result[:len(compressedPub)], compressedPub)
	copy(result[len(compressedPub):], symCiphertext)

	return base64.StdEncoding.EncodeToString(result), nil
}

func ecdhDeriveSharedSecret(privKey *ecdsa.PrivateKey, pubKey *ecdsa.PublicKey) ([]byte, error) {
	ecdhPriv, err := privKey.ECDH()
	if err != nil {
		return nil, fmt.Errorf("ecdsa private key ecdh failed: %s", err.Error())
	}

	ecdhPub, err := pubKey.ECDH()
	if err != nil {
		return nil, fmt.Errorf("ecdsa public key ecdh failed: %s", err.Error())
	}

	sharedSecret, err := ecdhPriv.ECDH(ecdhPub)
	if err != nil {
		return nil, fmt.Errorf("ecdh generate shared secret failed: %s", err.Error())
	}

	return sharedSecret, nil
}

func newAESGCMAndNONCE(sharedSecret []byte) (cipher.AEAD, []byte, error) {
	hkdfReader := hkdf.New(sha256.New, sharedSecret, nil, HKDF_INFO)
	key := make([]byte, AES_KEY_LENGTH)
	nonce := make([]byte, GCM_NONCE_LENGTH)
	if _, err := hkdfReader.Read(key); err != nil {
		return nil, nil, fmt.Errorf("hkdf read aes key failed: %s", err.Error())
	}

	if _, err := hkdfReader.Read(nonce); err != nil {
		return nil, nil, fmt.Errorf("hkdf read gcm nonce failed: %s", err.Error())
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("aes new cipher failed: %s", err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("aes new gcm failed: %s", err.Error())
	}

	return gcm, nonce, nil
}

func Decrypt(privKey *ecdsa.PrivateKey, encryptedData string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("base64 decode encrypted data failed: %s", err.Error())
	}

	if len(ciphertext) < COMPRESSED_PUB_LENGTH {
		return nil, fmt.Errorf("invalid ciphertext length %d", len(ciphertext))
	}

	curve := elliptic.P256()
	compressedPub := ciphertext[:COMPRESSED_PUB_LENGTH]
	symCiphertext := ciphertext[COMPRESSED_PUB_LENGTH:]
	x, y := elliptic.UnmarshalCompressed(curve, compressedPub)
	if x == nil {
		return nil, fmt.Errorf("invalid compressed public key")
	}

	sharedSecret, err := ecdhDeriveSharedSecret(privKey, &ecdsa.PublicKey{Curve: curve, X: x, Y: y})
	if err != nil {
		return nil, err
	}

	gcm, nonce, err := newAESGCMAndNONCE(sharedSecret)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, symCiphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("gcm decrypt ciphertext failed: %s", err.Error())
	}

	return plaintext, nil
}

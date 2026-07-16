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

var (
	HKDF_INFO = []byte("ecies-p256-aes128gcm")

	ecCurve = elliptic.P256()
)

type PrivateKey struct {
	*ecdsa.PrivateKey
}

type PublicKey struct {
	*ecdsa.PublicKey
}

func GenerateECIESKeyPair() (*PrivateKey, *PublicKey, error) {
	if privKey, err := ecdsa.GenerateKey(ecCurve, rand.Reader); err != nil {
		return nil, nil, err
	} else {
		return &PrivateKey{privKey}, &PublicKey{&privKey.PublicKey}, nil
	}
}

func (pub *PublicKey) WriteToFile(fileName string) error {
	pubKeyFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("create public key file failed: %w", err)
	}

	defer pubKeyFile.Close()

	pubKeyDER, err := x509.MarshalPKIXPublicKey(pub.PublicKey)
	if err != nil {
		return fmt.Errorf("x509 marshal public key failed: %w", err)
	}

	pubKeyPEM := &pem.Block{
		Type:  PKIX_PUBLIC_KEY_TYPE,
		Bytes: pubKeyDER,
	}
	if err := pem.Encode(pubKeyFile, pubKeyPEM); err != nil {
		return fmt.Errorf("write public key to file failed: %w", err)
	}

	return nil
}

func LoadPublicKeyFromFile(fileName string) (*PublicKey, error) {
	pubKeyPEM, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read public key failed: %w", err)
	}

	return ParsePublicKey(pubKeyPEM)
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

func ParsePublicKey(pubKeyPEM []byte) (*PublicKey, error) {
	if len(pubKeyPEM) == 0 {
		return nil, fmt.Errorf("public pem cannot be empty")
	}

	pubKeyBlock, _ := pem.Decode(pubKeyPEM)
	if pubKeyBlock == nil || pubKeyBlock.Type != PKIX_PUBLIC_KEY_TYPE {
		return nil, fmt.Errorf("cannot decode public key PEM")
	}

	publicKey, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509 parse public key failed: %w", err)
	}

	pubKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not ecdsa key")
	}

	return &PublicKey{pubKey}, nil
}

func (pub *PublicKey) Encrypt(plaintext []byte) (string, error) {
	if len(plaintext) == 0 {
		return "", fmt.Errorf("plaintext cannot be empty")
	}

	ephemeralPriv, err := ecdsa.GenerateKey(ecCurve, rand.Reader)
	if err != nil {
		return "", fmt.Errorf("ecdsa generate private key failed: %w", err)
	}

	sharedSecret, err := ecdhDeriveSharedSecret(ephemeralPriv, pub.PublicKey)
	if err != nil {
		return "", err
	}

	gcm, nonce, err := newAESGCMAndNONCE(sharedSecret)
	if err != nil {
		return "", err
	}

	symCiphertext := gcm.Seal(nil, nonce, plaintext, nil)
	compressedPub := elliptic.MarshalCompressed(ecCurve, ephemeralPriv.PublicKey.X, ephemeralPriv.PublicKey.Y)
	result := make([]byte, len(compressedPub)+len(symCiphertext))
	copy(result[:len(compressedPub)], compressedPub)
	copy(result[len(compressedPub):], symCiphertext)

	return base64.StdEncoding.EncodeToString(result), nil
}

func ecdhDeriveSharedSecret(privKey *ecdsa.PrivateKey, pubKey *ecdsa.PublicKey) ([]byte, error) {
	ecdhPriv, err := privKey.ECDH()
	if err != nil {
		return nil, fmt.Errorf("ecdsa private key ecdh failed: %w", err)
	}

	ecdhPub, err := pubKey.ECDH()
	if err != nil {
		return nil, fmt.Errorf("ecdsa public key ecdh failed: %w", err)
	}

	sharedSecret, err := ecdhPriv.ECDH(ecdhPub)
	if err != nil {
		return nil, fmt.Errorf("ecdh generate shared secret failed: %w", err)
	}

	return sharedSecret, nil
}

func newAESGCMAndNONCE(sharedSecret []byte) (cipher.AEAD, []byte, error) {
	hkdfReader := hkdf.New(sha256.New, sharedSecret, nil, HKDF_INFO)
	key := make([]byte, AES_KEY_LENGTH)
	nonce := make([]byte, GCM_NONCE_LENGTH)
	if _, err := hkdfReader.Read(key); err != nil {
		return nil, nil, fmt.Errorf("hkdf read aes key failed: %w", err)
	}

	if _, err := hkdfReader.Read(nonce); err != nil {
		return nil, nil, fmt.Errorf("hkdf read gcm nonce failed: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("aes new cipher failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("aes new gcm failed: %w", err)
	}

	return gcm, nonce, nil
}

func (priv *PrivateKey) Public() *PublicKey {
	return &PublicKey{&priv.PrivateKey.PublicKey}
}

func (priv *PrivateKey) WriteToFile(fileName string) error {
	privKeyFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("create private key file failed: %w", err)
	}

	defer privKeyFile.Close()

	privKeyDER, err := x509.MarshalPKCS8PrivateKey(priv.PrivateKey)
	if err != nil {
		return fmt.Errorf("x509 marshal private key failed: %w", err)
	}

	privKeyPEM := &pem.Block{
		Type:  PKCS8_PRIVATE_KEY_TYPE,
		Bytes: privKeyDER,
	}
	if err := pem.Encode(privKeyFile, privKeyPEM); err != nil {
		return fmt.Errorf("write private key to file failed: %w", err)
	}

	return nil
}

func LoadPrivateKeyFromFile(fileName string) (*PrivateKey, error) {
	privKeyPEM, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read private key failed: %w", err)
	}

	return ParsePrivateKey(privKeyPEM)
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

func ParsePrivateKey(privKeyPEM []byte) (*PrivateKey, error) {
	if len(privKeyPEM) == 0 {
		return nil, fmt.Errorf("private pem cannot be empty")
	}

	privKeyBlock, _ := pem.Decode(privKeyPEM)
	if privKeyBlock == nil || privKeyBlock.Type != PKCS8_PRIVATE_KEY_TYPE {
		return nil, fmt.Errorf("cannot decode private key PEM")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(privKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509 parse private key failed: %w", err)
	}

	privKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not ecdsa key")
	}

	return &PrivateKey{privKey}, nil
}

func (priv *PrivateKey) Decrypt(encryptedData string) ([]byte, error) {
	if len(encryptedData) == 0 {
		return nil, fmt.Errorf("ciphertext cannot be empty")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("base64 decode encrypted data failed: %w", err)
	}

	if len(ciphertext) < COMPRESSED_PUB_LENGTH {
		return nil, fmt.Errorf("invalid ciphertext length %d", len(ciphertext))
	}

	compressedPub := ciphertext[:COMPRESSED_PUB_LENGTH]
	symCiphertext := ciphertext[COMPRESSED_PUB_LENGTH:]
	x, y := elliptic.UnmarshalCompressed(ecCurve, compressedPub)
	if x == nil {
		return nil, fmt.Errorf("invalid compressed public key")
	}

	sharedSecret, err := ecdhDeriveSharedSecret(priv.PrivateKey, &ecdsa.PublicKey{Curve: ecCurve, X: x, Y: y})
	if err != nil {
		return nil, err
	}

	gcm, nonce, err := newAESGCMAndNONCE(sharedSecret)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, symCiphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("gcm decrypt ciphertext failed: %w", err)
	}

	return plaintext, nil
}

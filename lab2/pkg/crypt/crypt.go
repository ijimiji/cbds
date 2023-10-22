package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"

	"golang.org/x/crypto/ssh"
)

const bitSize = 2048

func GenerateRSAPair() ([]byte, []byte, error) {
	privateKeyRaw, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, nil, err
	}
	privDER := x509.MarshalPKCS1PrivateKey(privateKeyRaw)
	privateKey := pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	})

	publicKeyRaw, err := ssh.NewPublicKey(&privateKeyRaw.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	publicKey := publicKeyRaw.Marshal()

	return privateKey, publicKey, nil
}

func Encrypt(publicKeyBytes []byte, data []byte) ([]byte, error) {
	publicKeySSH, err := ssh.ParsePublicKey(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	parsedCryptoKey := publicKeySSH.(ssh.CryptoPublicKey)
	pubCrypto := parsedCryptoKey.CryptoPublicKey()
	publicKey := pubCrypto.(*rsa.PublicKey)

	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		data,
		nil)
	if err != nil {
		return nil, err
	}

	return encryptedBytes, nil
}

func Decrypt(privateKey []byte, data []byte) ([]byte, error) {
	// decrypted, err := base64.StdEncoding.DecodeString(string(data))
	// if err != nil {
	// 	return nil, err
	// }
	decrypted := data

	block, _ := pem.Decode(privateKey)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	res, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, decrypted, nil)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func EncryptAES(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func DecryptAES(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

type Cipher struct {
	key    []byte
	aesgcm cipher.AEAD
	nonce  []byte
}

func newCipher() (Cipher, error) {
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	return Cipher{key, aesgcm, nonce}, nil
}

func (c *Cipher) encrypt(plaintext []byte) ([]byte, error) {
	// Load your secret key from a safe place and reuse it across multiple
	// Seal/Open calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).

	// plaintext := []byte("exampleplaintext")
	// plaintext := []byte(msg)

	ciphertext := c.aesgcm.Seal(nil, c.nonce, plaintext, nil)
	fmt.Printf("%x\n", ciphertext)
	return ciphertext, nil
}

func (c *Cipher) decrypt(ciphertext []byte) ([]byte, error) {

	// Load your secret key from a safe place and reuse it across multiple
	// Seal/Open calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).
	// key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	// // ciphertext, _ := hex.DecodeString("c3aaa29f002ca75870806e44086700f62ce4d43e902b3888e23ceff797a7a471")
	// nonce, _ := hex.DecodeString("64a9433eae7ccceee2fc0eda")

	// block, err := aes.NewCipher(key)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// aesgcm, err := cipher.NewGCM(block)
	// if err != nil {
	// 	panic(err.Error())
	// }

	plaintext, err := c.aesgcm.Open(nil, c.nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%s\n", plaintext)
	return plaintext, nil
}

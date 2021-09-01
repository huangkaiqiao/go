package crypto

import (
	// "crypto/aes"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAes(t *testing.T) {
	plaintext := []byte("exampleplaintext")
	c, _ := newCipher()
	ciphertext, _ := c.encrypt(plaintext)
	decrypttext, _ := c.decrypt(ciphertext)
	assert.Equal(t, plaintext, decrypttext)
}

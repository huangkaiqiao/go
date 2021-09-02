package crypto

import (
	// "crypto/aes"

	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSha256FileList(t *testing.T) {
	testSha256File(t, "plaintext.txt", "6ea0f22a7539d6d22ba90d093462827884479561d9a83c8867a890499ba0c367")
}

func testSha256File(t *testing.T, path string, expected string) {
	result, _ := sha256File(path)
	assert.Equal(t, hex.EncodeToString(result), expected)
}

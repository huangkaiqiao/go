package mnemonic

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// {
//   "code": 20000,
//   "data": {
//     "btc_addr": "3PeAaq3zizMknZzANSuREmTdQXrBrjt7G6",
//     "eth_addr": "0x09d3ab4adacfd665b7576d010467a878f2015e56",
//     "fil_addr": "f1mw3yk2wwoicky6hb5takjhaxjmu7gwpndmjutna",
//     "mnemonic": "connect lazy leave tray fashion hurdle inquiry shiver mammal sting upper prison"
//   },
//   "msg": "操作成功",
//   "success": true
// }
func TestCryption(t *testing.T) {
	entropy := "connect lazy leave tray fashion hurdle inquiry shiver mammal sting upper prison"
	key := MakeKey(entropy)
	message := "I can't breathe. I can't make it，don't kill me. "
	ciphertext, _ := key.Encrypt([]byte(message))
	fmt.Printf("ciphertext=%s\n", hex.EncodeToString(ciphertext))
	plaintext, _ := key.Decrypt(ciphertext)
	assert.Equal(t, message, string(plaintext))
}

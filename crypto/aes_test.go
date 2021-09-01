package crypto

import (
	// "crypto/aes"

	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesGcmList(t *testing.T) {
	testAesGcm(t, "exampleplaintext")
	testAesGcm(t, "聊聊 Go 语言中的面向对象编程")
	testAesGcm(t, "我看到很多成年人在评论区冷嘲热讽, 等你们未成年的时候看谁帮你们说话")
}

func testAesGcm(t *testing.T, msg string) {
	plaintext := []byte(msg)
	c, _ := newCipher()
	ciphertext, _ := c.encrypt(plaintext)
	decrypttext, _ := c.decrypt(ciphertext)
	assert.Equal(t, plaintext, decrypttext)
}

func TestAesGcmFile(t *testing.T) {
	c, _ := newCipher()
	inpath := "plaintext.txt"
	outpath, _ := c.encryptFile(inpath)
	plainpath, _ := c.decryptFile(outpath)
	fmt.Println(plainpath)
	// TODO
	// assert.Equal(t, inpath, plainpath)
}

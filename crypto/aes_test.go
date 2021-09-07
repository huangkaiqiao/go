package crypto

import (
	// "crypto/aes"

	"os"
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
	c, _ := NewCipher(nil, nil)
	ciphertext, _ := c.Encrypt(plaintext)
	decrypttext, _ := c.Decrypt(ciphertext)
	assert.Equal(t, plaintext, decrypttext)
}

func TestAesGcmFileList(t *testing.T) {
	testAesGcmFile(t, "plaintext.txt")
	// testAesGcmFile(t, "gradle-7.2-bin.zip"
	// testAesGcmFile(t, "android-studio-2020.3.1.23-linux.tar.gz")
}

func testAesGcmFile(t *testing.T, inpath string) {
	c, _ := NewCipher(RandomKey(), nil)
	// inpath := "plaintext.txt"
	os.Remove(inpath + ".mn1")
	os.Remove(inpath + ".out")
	outpath, _ := c.EncryptFile(inpath, "")
	plainpath, _ := c.DecryptFile(outpath, "")
	// fmt.Println(plainpath)
	expected, _ := Sha256File(inpath)
	result, _ := Sha256File(plainpath)
	assert.Equal(t, expected, result)
}

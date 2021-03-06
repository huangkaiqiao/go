// 用于文档加密的包

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	// "math/rand"
)

// 加密文件的结构体,由密钥, aesgcm 对象, nonce 组成
type Cipher struct {
	key    []byte
	aesgcm cipher.AEAD
	nonce  []byte
}

// 生成随机 key 值
func RandomKey() []byte {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err.Error())
	}
	return key
}

// 通过 key 和 nonce 生成 Cipher 对象, 传 nil 返回随机值
func NewCipher(key []byte, nonce []byte) (Cipher, error) {
	// key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	if key == nil {
		key = RandomKey()
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	if nonce == nil {
		nonce = make([]byte, 12)
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			panic(err.Error())
		}
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("nonce:%d, size:%d", len(nonce), aesgcm.NonceSize())
	return Cipher{key, aesgcm, nonce}, nil
}

// 加密字节数组
func (c *Cipher) Encrypt(plaintext []byte) ([]byte, error) {
	ciphertext := c.aesgcm.Seal(nil, c.nonce, plaintext, nil)
	fmt.Printf("ciphertest=%x\n", ciphertext)
	return ciphertext, nil
}

const BUF_SIZE = 32 * 1024
const CIPHER_SIZE = BUF_SIZE + 16

// 加密文件, inpath 输入文件的路径, outpath 加密文件的路径(可以传空字符串), 默认返回 inpath + ".mn1" 结尾的路径
func (c *Cipher) EncryptFile(inpath string, outpath string) (string, error) {
	infile, err := os.Open(inpath)
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()
	fi, _ := infile.Stat()
	size := fi.Size()
	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	if outpath == "" {
		outpath = inpath + ".mn1"
	}
	outfile, err := os.OpenFile(outpath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()
	// The buffer size must be multiple of 16 bytes
	buf := make([]byte, BUF_SIZE)
	for {
		n := -1
		var err error
		if size < BUF_SIZE && 0 < size {
			buf1 := make([]byte, size)
			buf = make([]byte, BUF_SIZE)
			infile.Read(buf1)
			copy(buf, buf1)
			// fmt.Println(size)
			n, err = infile.Read(make([]byte, 1))
		} else {
			n, err = infile.Read(buf)
		}
		size = size - int64(n)
		if n > 0 {
			ciphertext := c.aesgcm.Seal(nil, c.nonce, buf, nil)
			// Write into file
			outfile.Write(ciphertext)
		}

		if err == io.EOF {
			ciphertext := c.aesgcm.Seal(nil, c.nonce, buf, nil)
			outfile.Write(ciphertext)
			break
		}

		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}
	// Append the IV
	outfile.Write(c.nonce)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(size))
	outfile.Write(b)
	return outfile.Name(), nil
}

// 获取加密后的文件的 nonce
func nonceOfFile(infile *os.File) ([]byte, error) {
	iv := make([]byte, 12)
	fi, err := infile.Stat()
	msgLen := fi.Size() - int64(len(iv)) - 8
	_, err = infile.ReadAt(iv, msgLen)
	if err != nil {
		log.Fatal(err)
	}
	return iv, nil
}

// 解密文件, inpath 加密文件路径, outpath 输出文件路径(可选), 默认返回 inpath + ".out"
func (c *Cipher) DecryptFile(inpath string, outpath string) (string, error) {
	infile, err := os.Open(inpath)
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	fi, err := infile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	// iv := make([]byte, 12)
	// msgLen := fi.Size() - int64(len(iv)) - 8
	// _, err = infile.ReadAt(iv, msgLen)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	iv, err := nonceOfFile(infile)
	if err != nil {
		log.Fatal(err)
	}
	c.nonce = iv // 更新了 cipher 的 nonce
	msgLen := fi.Size() - int64(len(iv)) - 8

	sizeOfLastBytes := make([]byte, 8)
	_, err = infile.ReadAt(sizeOfLastBytes, msgLen+12)
	sizeOfLast := int64(binary.LittleEndian.Uint64(sizeOfLastBytes))

	if err != nil {
		log.Fatal(err)
	}

	if outpath == "" {
		outpath = strings.TrimSuffix(inpath, filepath.Ext(inpath)) + ".out"
	}
	outfile, err := os.OpenFile(outpath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	// The buffer size must be multiple of 16 bytes
	buf := make([]byte, CIPHER_SIZE)
	for {
		n, err := infile.Read(buf)
		if n > 0 {
			// The last bytes are the IV, don't belong the original message
			msgLen -= int64(n)
			// stream.XORKeyStream(buf, buf[:n])
			plaintext, _ := c.aesgcm.Open(nil, c.nonce, buf, nil)
			// Write into file
			// fmt.Printf("plaintext=%v msgLen=%d\n", string(plaintext), msgLen)
			if msgLen > 0 {
				outfile.Write(plaintext)
			} else {
				outfile.Write(plaintext[0:sizeOfLast])
			}
		}

		if err == io.EOF || msgLen == 0 {
			break
		}

		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}
	return outpath, nil
}

// 解密字节数组
func (c *Cipher) Decrypt(ciphertext []byte) ([]byte, error) {

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

	// fmt.Printf("%s\n", plaintext)
	return plaintext, nil
}

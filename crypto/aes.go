package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	// "math/rand"
)

type Cipher struct {
	key    []byte
	aesgcm cipher.AEAD
	nonce  []byte
}

func randomKey() []byte {
	key := make([]byte, 32)
	return key
}

func newCipher() (Cipher, error) {
	// key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	key := randomKey()
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

	// nonce, _ := hex.DecodeString("64a9433eae7ccceee2fc0eda")

	ciphertext := c.aesgcm.Seal(nil, c.nonce, plaintext, nil)
	fmt.Printf("%x\n", ciphertext)
	return ciphertext, nil
}

func (c *Cipher) encryptFile(inpath string) (string, error) {
	infile, err := os.Open(inpath)
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()
	fi, _ := infile.Stat()
	size := fi.Size()
	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	outfile, err := os.OpenFile(inpath+".bin", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()
	// The buffer size must be multiple of 16 bytes
	buf := make([]byte, 1024)
	for {
		n := -1
		var err error
		if size < 1024 && 0 < size {
			buf1 := make([]byte, size)
			buf = make([]byte, 1024)
			n, err = infile.Read(buf1)
			copy(buf, buf1)
			fmt.Println(size)
		} else {
			n, err = infile.Read(buf)
		}
		size = size - int64(n)
		if n > 0 {
			//stream.XORKeyStream(buf, buf[:n])
			ciphertext := c.aesgcm.Seal(nil, c.nonce, buf, nil)
			// Write into file
			outfile.Write(ciphertext)
		}

		if err == io.EOF {
			fmt.Println("EOF")
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
	return outfile.Name(), nil
}

func (c *Cipher) decryptFile(inpath string) (string, error) {
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

	iv := make([]byte, 12)
	msgLen := fi.Size() - int64(len(iv))
	_, err = infile.ReadAt(iv, msgLen)
	if err != nil {
		log.Fatal(err)
	}

	outpath := strings.TrimSuffix(inpath, filepath.Ext(inpath)) + ".out"
	outfile, err := os.OpenFile(outpath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	// The buffer size must be multiple of 16 bytes
	buf := make([]byte, 1040)
	for {
		n, err := infile.Read(buf)
		if n > 0 {
			// The last bytes are the IV, don't belong the original message
			if n > int(msgLen) {
				n = int(msgLen)
			}
			msgLen -= int64(n)
			// stream.XORKeyStream(buf, buf[:n])
			plaintext, _ := c.aesgcm.Open(nil, c.nonce, buf, nil)
			// Write into file
			outfile.Write(plaintext)
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}
	return outpath, nil
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

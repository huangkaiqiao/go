package crypto

import (
	"crypto/sha256"
	"io"
	"log"
	"os"
)

func Sha256File(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%x", h.Sum(nil))
	return h.Sum(nil), nil
}

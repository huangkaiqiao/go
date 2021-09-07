# go-lib
golang crypto library

## Install(安装)
go get github.com/huangkaiqiao/go/crypto

## Test(测试)
```bash
pinkiepie@gummy:~/Git/pinkiepie/go$ make tests
....
coverage: 78.6% of statements
ok      github.com/huangkaiqiao/go-crypto/crypto        0.002s  coverage: 78.6% of statements
```

## Example(用例)

It's very unlikely, but possible, that a given index does not produce a valid 
private key. Error checking is skipped in this example for brevity but should be handled in real code. In such a case, a ErrInvalidPrivateKey is returned.

ErrInvalidPrivateKey should be handled by trying the next index for a child key.

Any valid private key will have a valid public key so that `Key.PublicKey()`
method never returns an error.

```go
package main

import (
  "github.com/huangkaiqiao/go/crypto"
  "fmt"
)

func main() {
    // 使用随机 key 生成 cipher 对象
    c, _ := NewCipher(RandomKey(), nil)
    inpath := "plaintext.txt"
    // os.Remove(inpath + ".mn1")
    // os.Remove(inpath + ".out")
    outpath, _ := c.EncryptFile(inpath, "")
    plainpath, _ := c.DecryptFile(outpath, "")
	// fmt.Println(plainpath)
	expected, _ := Sha256File(inpath)
	result, _ := Sha256File(plainpath)
	assert.Equal(t, expected, result)
}
```

## Thanks

The developers at [Factom](https://www.factom.com/) have contributed a lot to this library and have made many great improvements to it. Please check out their project(s) and give them a thanks if you use this library.

Thanks to [bartekn](https://github.com/bartekn) from Stellar for some important bug catches.


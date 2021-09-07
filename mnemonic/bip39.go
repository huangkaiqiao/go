package mnemonic

import (
	"github.com/filecoin-project/go-address"
	// "github.com/golang/glog"

	// "github.com/tyler-smith/go-bip32"
	// "github.com/huangkaiqiao/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type integer int

type number interface {
	bitSize() int
}

func (num integer) bitSize() int {
	return int(num)
}

func NewMnemonic(nums ...int) string {
	bitSize := 128
	for _, num := range nums {
		bitSize = integer(num).bitSize()
	}
	entropy, _ := bip39.NewEntropy(bitSize)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	// glog.Info("NewEntropy ", address)
	return string(mnemonic)
}

func FromMnemonic(entropy string) address.Address {

	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	// seed := bip39.NewSeed(mnemonic, "Secret Passphrase")
	seed := bip39.NewSeed(entropy, "")
	// masterKey, _ := bip32.NewMasterKey(seed)
	masterKey, _ := NewMasterKey(seed)
	// publicKey := masterKey.PublicKey()

	// Display mnemonic and keys
	// fmt.Println("Mnemonic: ", mnemonic)
	// fmt.Println("Master private key: ", masterKey)
	// fmt.Println("Master public key: ", publicKey)
	// m44h, _ := masterKey.NewChildKey(44 + bip32.FirstHardenedChild)
	// m44h_461h, _ := m44h.NewChildKey(461 + bip32.FirstHardenedChild)
	// m44h_461h_0h, _ := m44h_461h.NewChildKey(0 + bip32.FirstHardenedChild)
	m44h, _ := masterKey.NewChildKey(44 + FirstHardenedChild)
	m44h_461h, _ := m44h.NewChildKey(461 + FirstHardenedChild)
	m44h_461h_0h, _ := m44h_461h.NewChildKey(0 + FirstHardenedChild)
	m44h_461h_0h_0, _ := m44h_461h_0h.NewChildKey(0)
	m44h_461h_0h_0_0, _ := m44h_461h_0h_0.NewChildKey(0)
	pk := m44h_461h_0h_0_0
	pub := pk.UncompressedPublicKey()
	address.CurrentNetwork = address.Mainnet
	// address from a secp pub key
	secp256k1Address, _ := address.NewSecp256k1Address(pub.Key)
	return secp256k1Address
}

/* func main() {
	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, _ := bip39.NewEntropy(256)
	glog.Info("go-bip39.go ", entropy)
	// mnemonic, _ := bip39.NewMnemonic(entropy)
	mnemonic := "similar border inch plate inmate talent humor oil tuna hobby describe head"

	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	// seed := bip39.NewSeed(mnemonic, "Secret Passphrase")
	seed := bip39.NewSeed(mnemonic, "")

	masterKey, _ := bip32.NewMasterKey(seed)
	publicKey := masterKey.PublicKey()

	// Display mnemonic and keys
	fmt.Println("Mnemonic: ", mnemonic)
	fmt.Println("Master private key: ", masterKey)
	fmt.Println("Master public key: ", publicKey)
	m44h, _ := masterKey.NewChildKey(44 + bip32.FirstHardenedChild)
	m44h_461h, _ := m44h.NewChildKey(461 + bip32.FirstHardenedChild)
	m44h_461h_0h, _ := m44h_461h.NewChildKey(0 + bip32.FirstHardenedChild)
	m44h_461h_0h_0, _ := m44h_461h_0h.NewChildKey(0)
	m44h_461h_0h_0_0, _ := m44h_461h_0h_0.NewChildKey(0)
	aurora1 := m44h_461h_0h_0_0
	// fmt.Println("Filecoin private key: ", aurora1)
	glog.Info("Aurora1 private key: ", hex.EncodeToString(aurora1.Key))
	// fmt.Println("Filecoin base64: ", b64.StdEncoding.EncodeToString(aurora1.Key))
	aurora1Pub := aurora1.UncompressedPublicKey()
	glog.Info("Aurora1 public key: ", hex.EncodeToString(aurora1Pub.Key))
	// address from a secp pub key
	secp256k1Address, _ := address.NewSecp256k1Address(aurora1Pub.Key)
	glog.Info("go-bip39.go ", secp256k1Address)
	// zero := uint32(0x00000000)
	// m44h_0h, _ := m44h.NewChildKey(bip32.FirstHardenedChild)
	// m44h_0h_0h, _ := m44h_0h.NewChildKey(bip32.FirstHardenedChild)
	// m44h_0h_0h_0, _ := m44h_0h_0h.NewChildKey(0)
	// m44h_0h_0h_0_0, _ := m44h_0h_0h_0.NewChildKey(0)
	// account := m44h_0h_0h
	// bip32 := m44h_0h_0h_0
	// aurora1_btc := m44h_0h_0h_0_0
	// fmt.Println("Account Extended Private Key: ", account)
	// fmt.Println("BIP32 Extended Private Key: ", bip32)
	// fmt.Println("Bitcoin private key: ", aurora1_btc)
	// fmt.Println("Bitcoin base64: ", b64.StdEncoding.EncodeToString(aurora1_btc.Key))
}*/

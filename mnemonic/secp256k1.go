package mnemonic

import (
	ecies "github.com/ecies/go"
	"github.com/tyler-smith/go-bip39"
)

func MakeKey(entropy string) *Key {

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
	// pub := pk.UncompressedPublicKey()
	// address.CurrentNetwork = address.Mainnet
	// address from a secp pub key
	// secp256k1Address, _ := address.NewSecp256k1Address(pub.Key)
	// return secp256k1Address
	pk.EciesKey = *ecies.NewPrivateKeyFromBytes(pk.Key)
	return pk
}

func (key *Key) Encrypt(plaintext []byte) ([]byte, error) {
	return ecies.Encrypt(key.EciesKey.PublicKey, plaintext)
}

func (key *Key) Decrypt(ciphertext []byte) ([]byte, error) {
	return ecies.Decrypt(&key.EciesKey, ciphertext)
}

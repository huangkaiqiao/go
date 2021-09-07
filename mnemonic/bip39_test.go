// package lib
package mnemonic

import (
	"strings"
	"testing"
)

func TestNewMnemonic(t *testing.T) {
	var (
		expected1 = 12
		expected2 = 24
	)
	entropy := NewMnemonic()
	actual := len(strings.Split(entropy, " "))
	if actual != expected1 {
		t.Errorf("count of words in NewMnemonic() = %d; expected %d, entropy = %s", actual, expected1, entropy)
	}
	entropy = NewMnemonic(256)
	actual = len(strings.Split(entropy, " "))
	if actual != expected2 {
		t.Errorf("count of words in NewMnemonic() = %d; expected %d, entropy = %s", actual, expected2, entropy)
	}
}

func TestFromMnemonic(t *testing.T) {
	var (
		mnemonic = "similar border inch plate inmate talent humor oil tuna hobby describe head"
		expect   = "f1ftbz7dk462lneljmkkop22bmoafksmpzmzjeisa"
	)
	address := FromMnemonic(mnemonic)
	if address.String() != expect {
		t.Errorf("address of FromMnemonic() = %s, expected %s", address, expect)
	}
}

package main

import "testing"

func TestSignAndVerify(t *testing.T) {
	w := NewWallet()
	data := []byte("pay Bob 5 coins")
	sig := Sign(w.PrivateKey, data)

	if !Verify(w.PublicKey, data, sig) {
		t.Fatal("expected a genuine signature to verify")
	}
}

func TestVerifyFailsOnTamperedData(t *testing.T) {
	w := NewWallet()
	data := []byte("pay Bob 5 coins")
	sig := Sign(w.PrivateKey, data)

	tampered := []byte("pay Bob 9 coins")
	if Verify(w.PublicKey, tampered, sig) {
		t.Fatal("expected verification to fail for tampered data")
	}
}

func TestVerifyFailsWithWrongPublicKey(t *testing.T) {
	alice := NewWallet()
	mallory := NewWallet()
	data := []byte("pay Bob 5 coins")
	sig := Sign(alice.PrivateKey, data)

	if Verify(mallory.PublicKey, data, sig) {
		t.Fatal("expected verification to fail against an unrelated public key")
	}
}

func TestTwoWalletsDifferentAddresses(t *testing.T) {
	w1 := NewWallet()
	w2 := NewWallet()
	if w1.Address() == w2.Address() {
		t.Fatal("expected two independently generated wallets to have different addresses")
	}
}

func TestAddressIsDeterministicForSameKey(t *testing.T) {
	w := NewWallet()
	if w.Address() != w.Address() {
		t.Fatal("expected calling Address() twice on the same wallet to give the same result")
	}
}

func TestVerifyRejectsMalformedPublicKey(t *testing.T) {
	w := NewWallet()
	data := []byte("pay Bob 5 coins")
	sig := Sign(w.PrivateKey, data)

	truncated := w.PublicKey[:len(w.PublicKey)-10]
	if Verify(truncated, data, sig) {
		t.Fatal("expected Verify to reject a malformed/truncated public key")
	}
}

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"errors"
)

var curve = elliptic.P256()

func Hash(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte // uncompressed X||Y bytes
}

// marshalPublicKey turns an *ecdsa.PublicKey into raw X||Y bytes, each
// coordinate padded to the curve's byte width. (Provided for you.)
func marshalPublicKey(pub *ecdsa.PublicKey) []byte {
	byteLen := (pub.Curve.Params().BitSize + 7) / 8
	buf := make([]byte, 2*byteLen)
	pub.X.FillBytes(buf[:byteLen])
	pub.Y.FillBytes(buf[byteLen:])
	return buf
}

// unmarshalPublicKey reverses marshalPublicKey.
func unmarshalPublicKey(pubKey []byte) (*ecdsa.PublicKey, error) {
	byteLen := (curve.Params().BitSize + 7) / 8
	if len(pubKey) != 2*byteLen {
		return nil, errors.New("crypto: invalid public key length")
	}
	// TODO: reconstruct X and Y from pubKey using big.Int.SetBytes,
	// and return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}.
	panic("TODO: implement unmarshalPublicKey")
}

// NewWallet generates a fresh ECDSA key pair on the P-256 curve.
func NewWallet() *Wallet {
	panic("TODO: implement NewWallet using ecdsa.GenerateKey and crypto/rand")
}

// Address derives a short, printable identifier from the wallet's
// public key: hash the public key and hex-encode the result.
func (w *Wallet) Address() string {
	panic("TODO: implement Address")
}

// Sign signs the SHA-256 hash of data with priv, returning the signature.
func Sign(priv *ecdsa.PrivateKey, data []byte) []byte {
	panic("TODO: implement Sign using ecdsa.SignASN1")
}

// Verify reports whether signature is a valid signature over data's
// hash, produced by the private key matching pubKey.
func Verify(pubKey []byte, data, signature []byte) bool {
	panic("TODO: implement Verify using ecdsa.VerifyASN1")
}

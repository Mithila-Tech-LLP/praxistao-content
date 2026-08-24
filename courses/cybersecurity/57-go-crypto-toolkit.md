# Chapter 57: Go Cryptography Toolkit — Building Crypto Primitives

*Go's `crypto` package is one of the best standard libraries for cryptography. This chapter builds a complete toolkit: encryption, hashing, key derivation, signatures, and secure password storage.*

---

## Go's Crypto Packages

```
crypto/
├── aes           — AES symmetric encryption
├── cipher        — Block cipher modes (GCM, CTR, CBC)
├── ecdsa         — Elliptic Curve Digital Signatures
├── ed25519       — Ed25519 signatures
├── hmac          — HMAC message authentication
├── rand          — Cryptographically secure random numbers
├── rsa           — RSA asymmetric encryption
├── sha256        — SHA-256 hash
├── sha512        — SHA-512 hash
├── tls           — TLS 1.2/1.3
└── x509          — Certificate management

golang.org/x/crypto/
├── bcrypt        — Password hashing
├── argon2        — Memory-hard password hashing (recommended)
├── chacha20poly1305 — ChaCha20-Poly1305 AEAD
└── pbkdf2        — Key derivation function
```

---

## Secure Random Numbers

```go
package main

import (
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "math/big"
)

// Never use math/rand for security — it's predictable
// Always use crypto/rand

func generateToken(n int) string {
    b := make([]byte, n)
    if _, err := rand.Read(b); err != nil {
        panic(err)
    }
    return hex.EncodeToString(b)
}

func generateSecureInt(max int64) int64 {
    n, err := rand.Int(rand.Reader, big.NewInt(max))
    if err != nil {
        panic(err)
    }
    return n.Int64()
}

func main() {
    // 32-byte session token
    fmt.Println("Session token:", generateToken(32))
    
    // API key
    fmt.Println("API key:", generateToken(24))
    
    // Random number 0-99
    fmt.Println("Random:", generateSecureInt(100))
}
```

---

## Hashing

```go
package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "crypto/sha512"
    "encoding/hex"
    "fmt"
)

func hashSHA256(data []byte) string {
    h := sha256.Sum256(data)
    return hex.EncodeToString(h[:])
}

func hashSHA512(data []byte) string {
    h := sha512.Sum512(data)
    return hex.EncodeToString(h[:])
}

// HMAC — keyed hash (verifiable by holder of secret key)
// Use for: API request signing, message authentication
func hmacSHA256(data, key []byte) string {
    mac := hmac.New(sha256.New, key)
    mac.Write(data)
    return hex.EncodeToString(mac.Sum(nil))
}

func verifyHMAC(data, key []byte, expected string) bool {
    actual := hmacSHA256(data, key)
    // Use constant-time comparison to prevent timing attacks
    return hmac.Equal([]byte(actual), []byte(expected))
}

func main() {
    data := []byte("important message")
    key := []byte("secret-key-32-bytes-long-minimum!")
    
    fmt.Printf("SHA-256: %s\n", hashSHA256(data))
    fmt.Printf("HMAC:    %s\n", hmacSHA256(data, key))
    
    // Verify
    mac := hmacSHA256(data, key)
    fmt.Printf("Valid:   %v\n", verifyHMAC(data, key, mac))
    
    // Tampered data
    fmt.Printf("Tampered valid: %v\n", verifyHMAC([]byte("tampered"), key, mac))
}
```

---

## Symmetric Encryption (AES-256-GCM)

AES-GCM is the right choice: authenticated encryption (AEAD) — provides both confidentiality AND integrity.

```go
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "errors"
    "io"
)

// Encrypt encrypts plaintext using AES-256-GCM
// Returns: nonce (12 bytes) + ciphertext + tag (16 bytes)
func Encrypt(plaintext, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    // Random nonce — NEVER reuse a nonce with the same key
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    
    // Seal appends ciphertext + tag to nonce
    return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts data encrypted with Encrypt()
func Decrypt(ciphertext, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }
    
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}

// GenerateKey generates a 32-byte (256-bit) AES key
func GenerateKey() ([]byte, error) {
    key := make([]byte, 32)
    _, err := rand.Read(key)
    return key, err
}
```

---

## Password Hashing (Argon2)

```go
package main

import (
    "crypto/rand"
    "crypto/subtle"
    "encoding/hex"
    "fmt"
    
    "golang.org/x/crypto/argon2"
)

type HashParams struct {
    Memory      uint32
    Iterations  uint32
    Parallelism uint8
    SaltLen     uint32
    KeyLen      uint32
}

// Recommended parameters (2024)
var defaultParams = &HashParams{
    Memory:      64 * 1024,  // 64 MB
    Iterations:  3,
    Parallelism: 2,
    SaltLen:     16,
    KeyLen:      32,
}

func HashPassword(password string) (string, error) {
    salt := make([]byte, defaultParams.SaltLen)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }
    
    hash := argon2.IDKey(
        []byte(password),
        salt,
        defaultParams.Iterations,
        defaultParams.Memory,
        defaultParams.Parallelism,
        defaultParams.KeyLen,
    )
    
    // Store as "salt$hash"
    return hex.EncodeToString(salt) + "$" + hex.EncodeToString(hash), nil
}

func VerifyPassword(password, stored string) bool {
    parts := split(stored, "$")
    if len(parts) != 2 {
        return false
    }
    
    salt, _ := hex.DecodeString(parts[0])
    storedHash, _ := hex.DecodeString(parts[1])
    
    hash := argon2.IDKey(
        []byte(password),
        salt,
        defaultParams.Iterations,
        defaultParams.Memory,
        defaultParams.Parallelism,
        defaultParams.KeyLen,
    )
    
    // Constant-time comparison — prevent timing attacks
    return subtle.ConstantTimeCompare(hash, storedHash) == 1
}

func split(s, sep string) []string {
    idx := len(s)
    for i, r := range s {
        if string(r) == sep {
            idx = i
            break
        }
    }
    if idx == len(s) {
        return []string{s}
    }
    return []string{s[:idx], s[idx+1:]}
}

func main() {
    hash, _ := HashPassword("MySecurePassword123!")
    fmt.Println("Hash:", hash)
    
    fmt.Println("Correct:", VerifyPassword("MySecurePassword123!", hash))
    fmt.Println("Wrong:  ", VerifyPassword("wrongpassword", hash))
}
```

---

## Digital Signatures (Ed25519)

```go
package main

import (
    "crypto/ed25519"
    "crypto/rand"
    "encoding/hex"
    "fmt"
)

func generateSigningKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
    return ed25519.GenerateKey(rand.Reader)
}

func signData(data []byte, privateKey ed25519.PrivateKey) []byte {
    return ed25519.Sign(privateKey, data)
}

func verifySignature(data, signature []byte, publicKey ed25519.PublicKey) bool {
    return ed25519.Verify(publicKey, data, signature)
}

func main() {
    pub, priv, err := generateSigningKeyPair()
    if err != nil {
        panic(err)
    }
    
    message := []byte("This message is authentic")
    sig := signData(message, priv)
    
    fmt.Printf("Message:   %s\n", message)
    fmt.Printf("Public key: %s\n", hex.EncodeToString(pub))
    fmt.Printf("Signature:  %s...\n", hex.EncodeToString(sig[:16]))
    fmt.Printf("Valid:      %v\n", verifySignature(message, sig, pub))
    
    // Tampered message
    tampered := []byte("This message is fake")
    fmt.Printf("Tampered:   %v\n", verifySignature(tampered, sig, pub))
}
```

---

## Key Derivation (PBKDF2)

```go
import (
    "crypto/sha256"
    "golang.org/x/crypto/pbkdf2"
)

// Derive an AES key from a password (for file encryption)
func deriveKey(password string, salt []byte) []byte {
    return pbkdf2.Key(
        []byte(password),
        salt,
        600_000,  // iterations (NIST 2024 recommendation)
        32,       // 256-bit key
        sha256.New,
    )
}
```

---

## Complete Crypto Toolkit Usage

```go
// Putting it all together: encrypted, signed file
package main

import (
    "fmt"
    "encoding/hex"
    "crypto/rand"
    "crypto/ed25519"
)

func encryptAndSign(plaintext []byte, encKey []byte, signingKey ed25519.PrivateKey) ([]byte, []byte, error) {
    // Encrypt
    ciphertext, err := Encrypt(plaintext, encKey)
    if err != nil {
        return nil, nil, err
    }
    
    // Sign the ciphertext (encrypt-then-sign)
    sig := ed25519.Sign(signingKey, ciphertext)
    
    return ciphertext, sig, nil
}

func decryptAndVerify(ciphertext, sig []byte, encKey []byte, pubKey ed25519.PublicKey) ([]byte, error) {
    // Verify signature first
    if !ed25519.Verify(pubKey, ciphertext, sig) {
        return nil, fmt.Errorf("signature verification failed")
    }
    
    // Decrypt
    return Decrypt(ciphertext, encKey)
}

func main() {
    key, _ := GenerateKey()
    pub, priv, _ := generateSigningKeyPair()
    
    message := []byte("Secret data for GoShield")
    
    ct, sig, _ := encryptAndSign(message, key, priv)
    fmt.Printf("Encrypted (%d bytes), signed (%d bytes)\n", len(ct), len(sig))
    
    decrypted, err := decryptAndVerify(ct, sig, key, pub)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Decrypted: %s\n", decrypted)
    _ = hex.EncodeToString(pub)
    _ = rand.Reader
}
```

---

## Summary

| Use Case | Algorithm | Go Package |
|----------|-----------|-----------|
| Symmetric encryption | AES-256-GCM | `crypto/aes`, `crypto/cipher` |
| Password hashing | Argon2id | `golang.org/x/crypto/argon2` |
| Digital signatures | Ed25519 | `crypto/ed25519` |
| Message authentication | HMAC-SHA256 | `crypto/hmac`, `crypto/sha256` |
| Key derivation | PBKDF2 / Argon2 | `golang.org/x/crypto/pbkdf2` |
| Secure random | CSPRNG | `crypto/rand` |

---

## Exercises

1. Build a file encryption tool in Go: encrypt with AES-256-GCM, sign with Ed25519, store as JSON
2. Implement a secure API request signing system: HMAC-sign each request with a shared secret
3. Create a password manager skeleton: store bcrypt/argon2 hashes, verify correctly
4. Research "cryptographic agility" — why is it both useful and dangerous in protocol design?

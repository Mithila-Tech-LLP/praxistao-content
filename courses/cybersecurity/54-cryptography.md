# Chapter 54: Cryptography — The Foundation of Security

*Cryptography isn't just TLS and passwords. It's the mathematical backbone of almost every security primitive. Understanding it lets you know when it's correctly applied and — more importantly — when it's dangerously broken.*

---

## Why Cryptography Matters for Security Professionals

- Understand why certain attacks work (weak keys, broken algorithms)
- Know when code is doing crypto wrong (ECB mode, MD5 passwords)
- Build tools that correctly handle sensitive data
- Read vulnerability disclosures that discuss crypto weaknesses
- Pass interviews that always ask about hashing vs encryption

---

## Foundational Concepts

### The Three Goals of Cryptography

| Goal | Question | Tool |
|------|---------|------|
| **Confidentiality** | Can only intended parties read it? | Encryption |
| **Integrity** | Was it tampered with? | Hashing, MACs |
| **Authenticity** | Did it really come from who I think? | Digital signatures |

### Kerckhoffs's Principle

*"A cryptosystem should be secure even if everything about the system, except the key, is public knowledge."*

Good crypto doesn't rely on secret algorithms — only secret keys. Never roll your own crypto (RYOC).

---

## Symmetric Encryption

Same key encrypts and decrypts. Fast. Scales well.

```
Plaintext + Key → [Encrypt] → Ciphertext
Ciphertext + Key → [Decrypt] → Plaintext
```

### AES (Advanced Encryption Standard)

The gold standard for symmetric encryption. Block cipher: operates on 128-bit blocks.

**Key sizes:**
- AES-128: 128-bit key (16 bytes) — still secure
- AES-192: 192-bit key (24 bytes)
- AES-256: 256-bit key (32 bytes) — preferred for sensitive data

### Block Cipher Modes

How AES handles data larger than 128 bits:

**ECB (Electronic Codebook) — NEVER USE:**
```
Block 1 → AES_encrypt(Block 1) = Cipher 1
Block 2 → AES_encrypt(Block 2) = Cipher 2
```
Identical plaintext blocks produce identical ciphertext blocks. Reveals patterns.

The ECB penguin problem: encrypt an image with ECB and you can still see the original shape.

**CBC (Cipher Block Chaining) — Old standard:**
```
Block 1 XOR IV → AES_encrypt → Cipher 1
Block 2 XOR Cipher 1 → AES_encrypt → Cipher 2
```
Each block depends on the previous. Needs random IV (initialization vector).
Vulnerable to **padding oracle attacks** (POODLE).

**GCM (Galois/Counter Mode) — Use this:**
```
Stream cipher mode + authentication tag
AES-256-GCM: encryption + integrity in one operation
```
- Fast
- Parallelizable
- Provides authenticated encryption (AEAD)
- **Nonce must NEVER be reused** — reusing nonce with same key breaks security entirely

### Go: AES-GCM Encryption

```go
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "errors"
    "fmt"
    "io"
)

func encrypt(key, plaintext []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    // Nonce must be unique per encryption with same key
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    // Seal appends ciphertext+tag to nonce
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return ciphertext, nil
}

func decrypt(key, ciphertext []byte) ([]byte, error) {
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
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err // authentication failed = tampered!
    }
    return plaintext, nil
}

func main() {
    // Generate a 256-bit key (never hardcode in real code!)
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        panic(err)
    }

    msg := []byte("This is a secret message from the security tool")
    
    encrypted, err := encrypt(key, msg)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Encrypted (%d bytes): %x...\n", len(encrypted), encrypted[:16])

    decrypted, err := decrypt(key, encrypted)
    if err != nil {
        panic("Decryption failed: " + err.Error())
    }
    fmt.Printf("Decrypted: %s\n", decrypted)
    
    // Tamper test: modify one byte
    encrypted[20] ^= 0xFF
    _, err = decrypt(key, encrypted)
    fmt.Printf("Tampered decryption error: %v\n", err)
    // Output: authentication tag mismatch — GCM detected tampering!
}
```

---

## Asymmetric Encryption

Different keys for encryption and decryption. Slower. Solves the key distribution problem.

```
Public key (shareable) + Private key (secret)

Encryption:   Sender encrypts with recipient's PUBLIC key
              Only recipient's PRIVATE key can decrypt

Signing:      Signer signs with their own PRIVATE key
              Anyone can verify with signer's PUBLIC key
```

**RSA:** Older, based on factoring large numbers. Still used in TLS, SSH.
- RSA-2048: minimum for new deployments
- RSA-4096: for long-lived certificates

**ECDSA/EdDSA:** Elliptic curve. Smaller keys, faster, better security per bit.
- P-256 (NIST): widely used
- **Ed25519:** Best choice for new systems — fast, secure, small keys

### Go: RSA Signing and Verification

```go
package main

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "fmt"
)

func main() {
    // Generate RSA key pair
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        panic(err)
    }
    publicKey := &privateKey.PublicKey

    // Message to sign
    message := []byte("This config was approved by security team")
    
    // Sign: hash the message, sign the hash
    hash := sha256.Sum256(message)
    signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
    if err != nil {
        panic(err)
    }
    fmt.Printf("Signature (%d bytes): %x...\n", len(signature), signature[:16])

    // Verify: check signature matches
    err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
    if err != nil {
        fmt.Println("Invalid signature!")
    } else {
        fmt.Println("Signature valid — message authentic")
    }

    // Tamper test
    message[0] = 'X'
    hash2 := sha256.Sum256(message)
    err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash2[:], signature)
    fmt.Printf("Tampered verification: %v\n", err)
    // verification error
}
```

---

## Hashing

One-way functions: easy to compute, infeasible to reverse.

**Properties:**
1. **Deterministic:** Same input → same output always
2. **One-way:** Cannot derive input from output
3. **Collision-resistant:** Different inputs should produce different outputs
4. **Avalanche effect:** Tiny change in input → completely different output

### Hash Algorithms

| Algorithm | Output | Status |
|-----------|--------|--------|
| **MD5** | 128-bit | BROKEN — never use for security |
| **SHA-1** | 160-bit | BROKEN for collision resistance |
| **SHA-256** | 256-bit | Secure — file integrity, TLS |
| **SHA-512** | 512-bit | Secure — higher security margin |
| **bcrypt** | Variable | Designed for passwords — SLOW by design |
| **Argon2** | Variable | Password hashing winner — SLOW, memory-hard |

### Why MD5/SHA-1 Are Broken for Passwords

```
MD5("password") = 5f4dcc3b5aa765d61d8327deb882cf99

Two problems:
1. Fast: GPU can compute 10 billion MD5 hashes/second
2. Rainbow tables exist: precomputed hash→plaintext databases

Cracking "password" from its MD5 hash takes milliseconds.
```

### Password Hashing

**Never store plaintext passwords. Never use MD5/SHA-1 for passwords.**

```go
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
    // Cost 12 = ~250ms to hash → slow enough to deter brute force
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(hash), err
}

func checkPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func main() {
    password := "secretpassword123"
    
    hash, err := hashPassword(password)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Hash: %s\n", hash)
    // $2a$12$... — bcrypt includes salt in the output
    
    fmt.Println("Correct:", checkPassword("secretpassword123", hash))  // true
    fmt.Println("Wrong:", checkPassword("wrongpassword", hash))        // false
}
```

### HMAC — Keyed Hashing

MAC = Message Authentication Code. Proves data came from someone with the key AND wasn't tampered with.

```go
package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "fmt"
)

func createHMAC(key, message []byte) []byte {
    mac := hmac.New(sha256.New, key)
    mac.Write(message)
    return mac.Sum(nil)
}

func verifyHMAC(key, message, signature []byte) bool {
    expected := createHMAC(key, message)
    return hmac.Equal(expected, signature)
    // hmac.Equal is constant-time: prevents timing attacks
}

func main() {
    key := []byte("super-secret-key-32-bytes-long!!")
    msg := []byte("log entry: user 42 deleted file config.json")
    
    sig := createHMAC(key, msg)
    fmt.Printf("HMAC: %x\n", sig)
    
    fmt.Println("Valid:", verifyHMAC(key, msg, sig))
    
    // Tamper with message
    msg[5] = 'X'
    fmt.Println("Tampered:", verifyHMAC(key, msg, sig))  // false
}
```

---

## Key Exchange — Diffie-Hellman

The "magic" that lets two parties establish a shared secret over an insecure channel.

```
1. Alice and Bob agree on public parameters (g=2, p=large prime)
2. Alice picks secret a=6, sends: A = g^a mod p = 2^6 mod p
3. Bob picks secret b=15, sends: B = g^b mod p = 2^15 mod p
4. Alice computes: s = B^a mod p = (g^b)^a mod p
5. Bob computes:   s = A^b mod p = (g^a)^b mod p
6. Both arrive at s = g^(ab) mod p — the shared secret!
7. Attacker sees: g, p, A, B — cannot derive s without solving discrete log
```

Modern variant: **ECDH (Elliptic Curve Diffie-Hellman)** — smaller, faster, same security.

**Forward secrecy:** If long-term key is compromised later, past sessions stay secret.
- Without FS: Attacker records encrypted traffic, later cracks the key, decrypts everything
- With FS: Each session uses fresh ECDH keys, past sessions cannot be decrypted

---

## TLS in Practice

TLS 1.3 (current standard) handshake:

```
Client → ClientHello (supported ciphers, key_share)
Server → ServerHello (chosen cipher, key_share, certificate)
         [Encrypted data starts HERE — 1-RTT only]
Client → [Verify certificate] → Finished
Server → Finished
```

**Cipher suite example:** `TLS_AES_256_GCM_SHA384`
- `AES_256_GCM` = symmetric encryption after handshake
- `SHA384` = PRF (pseudorandom function) for key derivation

### Checking TLS Configuration

```bash
# Check supported TLS versions and ciphers
openssl s_client -connect example.com:443 2>/dev/null | head -30

# Full TLS audit
# testssl.sh is the gold standard
./testssl.sh https://example.com

# Python: check cert expiry
echo | openssl s_client -connect example.com:443 2>/dev/null | \
    openssl x509 -noout -dates
```

### Go: TLS Client

```go
package main

import (
    "crypto/tls"
    "fmt"
    "net/http"
)

func main() {
    // Secure client (verify certificates)
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                MinVersion: tls.VersionTLS12,
                // Don't set InsecureSkipVerify: true in production!
            },
        },
    }

    resp, err := client.Get("https://example.com")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    // Check the TLS info
    tlsInfo := resp.TLS
    fmt.Printf("TLS Version: %x\n", tlsInfo.Version)
    fmt.Printf("Cipher Suite: %s\n", tls.CipherSuiteName(tlsInfo.CipherSuite))
    
    if len(tlsInfo.PeerCertificates) > 0 {
        cert := tlsInfo.PeerCertificates[0]
        fmt.Printf("Certificate: %s\n", cert.Subject.CommonName)
        fmt.Printf("Expires: %s\n", cert.NotAfter)
        fmt.Printf("Issuer: %s\n", cert.Issuer.Organization)
    }
}
```

---

## Common Crypto Mistakes (Attack Targets)

| Mistake | Problem | Attack |
|---------|---------|--------|
| MD5/SHA1 passwords | Fast to crack | Password cracking |
| ECB mode | Pattern leakage | Block structure analysis |
| Reused nonce/IV | Key recovery | Nonce reuse attack |
| Weak RNG | Predictable keys | Key prediction |
| `alg: none` JWT | No signature | Token forgery |
| Short RSA key (<2048) | Factorable | Key recovery |
| Self-signed cert, no validation | MitM possible | SSL stripping |
| Hardcoded keys in source | Key exposure | Config scraping |

---

## Summary

| Concept | Use case | Go package |
|---------|---------|-----------|
| AES-GCM | Data encryption | `crypto/aes`, `crypto/cipher` |
| RSA/ECDSA | Signatures, TLS | `crypto/rsa`, `crypto/ecdsa` |
| SHA-256 | File hashing, integrity | `crypto/sha256` |
| bcrypt | Password storage | `golang.org/x/crypto/bcrypt` |
| HMAC | Message authentication | `crypto/hmac` |
| TLS | Transport security | `crypto/tls` |
| rand.Reader | Cryptographic randomness | `crypto/rand` |

**Golden rules:**
1. Use `crypto/rand` not `math/rand` for any security-sensitive randomness
2. Never roll your own crypto — use established libraries
3. Passwords: bcrypt or argon2, cost factor appropriate for your hardware
4. Encryption: AES-256-GCM with random nonce
5. Signatures: Ed25519 for new systems

---

## Exercises

1. Implement file encryption using AES-256-GCM in Go. Encrypt `/etc/hosts` and verify you can decrypt it correctly.
2. Generate an Ed25519 key pair. Sign a message. Verify the signature.
3. Implement a simple password manager that stores passwords encrypted with a master key derived from a passphrase using `golang.org/x/crypto/scrypt`.
4. Use `openssl` to inspect the TLS certificate of 3 websites. What algorithms do they use?
5. Find a real-world example of an ECB-mode vulnerability or a reused-nonce vulnerability. Explain how it was exploited.

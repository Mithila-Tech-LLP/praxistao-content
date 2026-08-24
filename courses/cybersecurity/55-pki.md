# Chapter 55: PKI — Public Key Infrastructure

*PKI is the system that makes trust on the internet possible. It's why you trust your bank's website, why HTTPS works, and why code signing prevents malware from masquerading as legitimate software.*

---

## The Trust Problem

How do you know you're talking to the real google.com and not an impersonator?

Without PKI:
```
You → "I want to connect to google.com" → Network → Impersonator
                                                        ↓
                                              Impersonator sends you
                                              "I am Google!" (with no way to verify)
```

With PKI:
```
You → Connect to google.com → Google sends certificate
Certificate: "I am google.com. Signed by Google Trust Services.
              And Google Trust Services is signed by a Root CA you already trust."
Browser verifies signature chain → TRUST ESTABLISHED
```

---

## PKI Components

```
Root CA (Certificate Authority)
├── Self-signed — trusted by browsers/OS
├── Examples: DigiCert, Let's Encrypt, Comodo
└── Certificate stored in browser/OS trust store

Intermediate CA
├── Signed by Root CA
├── Issues certificates to websites
└── Compromise of intermediate doesn't compromise Root CA

End-Entity Certificate (leaf)
├── Issued to a specific domain/server
├── Contains: domain name, public key, validity period
└── Signed by Intermediate CA

Chain of Trust:
Root CA → signs → Intermediate CA → signs → Server Certificate
```

---

## X.509 Certificate Structure

```
Version: 3
Serial Number: 0x2a4f...
Algorithm: sha256WithRSAEncryption
Issuer: C=US, O=Let's Encrypt, CN=R3
Validity:
    Not Before: Jan 1 2025
    Not After: Mar 31 2025    ← 90 day certs are standard
Subject: CN=www.example.com
Subject Public Key Info:
    Algorithm: rsaEncryption
    Public Key: (2048 bit)
X509v3 extensions:
    Subject Alt Name: DNS:example.com, DNS:www.example.com
    Key Usage: Digital Signature, Key Encipherment
    Extended Key Usage: TLS Web Server Authentication
    Basic Constraints: CA:FALSE
Signature: (signed by Issuer's private key)
```

---

## Working with Certificates

```bash
# Inspect a certificate
openssl s_client -connect google.com:443 2>/dev/null | openssl x509 -text -noout

# Decode a certificate file
openssl x509 -in cert.pem -text -noout

# Check certificate expiry
openssl x509 -in cert.pem -noout -enddate
# Output: notAfter=Mar 31 12:00:00 2025 GMT

# Verify certificate chain
openssl verify -CAfile /etc/ssl/certs/ca-certificates.crt cert.pem

# Generate a self-signed certificate (for dev/testing)
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 365 \
    -subj "/CN=localhost" -nodes

# Generate CSR (Certificate Signing Request — to send to CA)
openssl req -newkey rsa:2048 -keyout private.key -out request.csr \
    -subj "/CN=example.com/O=My Company"

# Let's Encrypt — free certificates
certbot --nginx -d example.com -d www.example.com
```

---

## Certificate Pinning

Normal TLS: trust any cert signed by any trusted CA.
Certificate pinning: trust ONLY this specific cert/public key.

```go
// Go: certificate pinning
import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "net/http"
)

func pinnedHTTPClient(expectedPubKeyHash string) *http.Client {
    return &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                VerifyPeerCertificate: func(rawCerts [][]byte, chains [][]*x509.Certificate) error {
                    if len(rawCerts) == 0 {
                        return fmt.Errorf("no certificate")
                    }
                    
                    cert, err := x509.ParseCertificate(rawCerts[0])
                    if err != nil {
                        return err
                    }
                    
                    // Hash the public key
                    pubKeyHash := fmt.Sprintf("%x", cert.RawSubjectPublicKeyInfo)
                    
                    if pubKeyHash != expectedPubKeyHash {
                        return fmt.Errorf("certificate pin mismatch: got %s, want %s",
                            pubKeyHash[:16], expectedPubKeyHash[:16])
                    }
                    return nil
                },
            },
        },
    }
}
```

---

## Private CA with Go

```go
package main

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "math/big"
    "os"
    "time"
)

func generateCA() (*rsa.PrivateKey, *x509.Certificate, error) {
    // Generate CA private key
    caKey, err := rsa.GenerateKey(rand.Reader, 4096)
    if err != nil {
        return nil, nil, err
    }
    
    // CA certificate template
    template := &x509.Certificate{
        SerialNumber: big.NewInt(1),
        Subject: pkix.Name{
            Organization: []string{"My Lab CA"},
            CommonName:   "My Lab Root CA",
        },
        NotBefore:             time.Now(),
        NotAfter:              time.Now().AddDate(10, 0, 0),
        IsCA:                  true,
        KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
        BasicConstraintsValid: true,
    }
    
    // Self-sign
    certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &caKey.PublicKey, caKey)
    if err != nil {
        return nil, nil, err
    }
    
    caCert, _ := x509.ParseCertificate(certBytes)
    
    // Write CA cert to file
    f, _ := os.Create("ca.crt")
    pem.Encode(f, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
    f.Close()
    
    return caKey, caCert, nil
}

func issueServerCert(caKey *rsa.PrivateKey, caCert *x509.Certificate, domain string) error {
    serverKey, _ := rsa.GenerateKey(rand.Reader, 2048)
    
    template := &x509.Certificate{
        SerialNumber: big.NewInt(2),
        Subject: pkix.Name{
            CommonName: domain,
        },
        DNSNames: []string{domain},
        NotBefore: time.Now(),
        NotAfter:  time.Now().AddDate(1, 0, 0),
        KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
        ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
    }
    
    certBytes, err := x509.CreateCertificate(rand.Reader, template, caCert, &serverKey.PublicKey, caKey)
    if err != nil {
        return err
    }
    
    // Write cert and key
    certFile, _ := os.Create(domain + ".crt")
    pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
    certFile.Close()
    
    keyFile, _ := os.Create(domain + ".key")
    pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(serverKey)})
    keyFile.Close()
    
    return nil
}

func main() {
    caKey, caCert, err := generateCA()
    if err != nil {
        panic(err)
    }
    
    if err := issueServerCert(caKey, caCert, "internal.company.com"); err != nil {
        panic(err)
    }
    
    println("CA and server certificate generated!")
    println("Add ca.crt to trust store, use internal.company.com.crt + .key in your server")
}
```

---

## PKI Attacks

```
Certificate Authority Compromise
├── DigiNotar (2011) — Iranian hackers, Google certs forged
├── SSL stripping — downgrade HTTPS to HTTP
└── Rogue CA — install malicious root cert

BGP Hijacking + Cert Issuance
├── Hijack IP route for domain.com
├── Get certificate from CA (domain validation uses HTTP)
└── Man-in-the-middle with valid certificate

CT Log Monitoring (defense)
├── All public certs must be logged in Certificate Transparency
├── Monitor crt.sh for unexpected certs for your domain
└── HPKP (deprecated) / CAA records restrict which CAs can issue
```

---

## Summary

| Concept | Purpose |
|---------|---------|
| Root CA | Top-level trust anchor (pre-installed in browser) |
| Intermediate CA | Issues end-entity certs, limits Root CA risk |
| X.509 certificate | Binds a public key to a domain name |
| Certificate chain | Root → Intermediate → Server cert |
| Let's Encrypt | Free, automated, 90-day certificates |
| Certificate pinning | Trust only specific cert, not any CA-signed cert |

---

## Exercises

1. Generate a self-signed CA and use it to issue a server cert for localhost — configure nginx to use it
2. Use `openssl s_client` to inspect the certificate chain of 3 different websites
3. Set up monitoring with crt.sh — watch for new certificates issued for a domain you own
4. Implement the Go private CA and issue a certificate — verify it with openssl verify

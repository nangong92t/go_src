package jwt

import (
"fmt"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

// Implements the RSA family of signing methods signing methods
type SigningMethodRSA struct {
	Name string
	Hash crypto.Hash
}

// Specific instances for RS256 and company
var (
	SigningMethodRS256 *SigningMethodRSA
	SigningMethodRS384 *SigningMethodRSA
	SigningMethodRS512 *SigningMethodRSA
)

func init() {
	// RS256
	SigningMethodRS256 = &SigningMethodRSA{"RS256", crypto.SHA256}
	RegisterSigningMethod(SigningMethodRS256.Alg(), func() SigningMethod {
		return SigningMethodRS256
	})

	// RS384
	SigningMethodRS384 = &SigningMethodRSA{"RS384", crypto.SHA384}
	RegisterSigningMethod(SigningMethodRS384.Alg(), func() SigningMethod {
		return SigningMethodRS384
	})

	// RS512
	SigningMethodRS512 = &SigningMethodRSA{"RS512", crypto.SHA512}
	RegisterSigningMethod(SigningMethodRS512.Alg(), func() SigningMethod {
		return SigningMethodRS512
	})
}

func (m *SigningMethodRSA) Alg() string {
	return m.Name
}

// Implements the Verify method from SigningMethod
// For this signing method, must be either a PEM encoded PKCS1 or PKCS8 RSA public key as
// []byte, or an rsa.PublicKey structure.
func (m *SigningMethodRSA) Verify(signingString, signature string, key interface{}) error {
	var err error

	// Decode the signature
	var sig []byte
fmt.Println("right: ")
	if sig, err = DecodeSegment(signature); err != nil {
		return err
	}

	var rsaKey *rsa.PublicKey

	switch k := key.(type) {
	case []byte:
		if rsaKey, err = ParseRSAPublicKeyFromPEM(k); err != nil {
			return err
		}
	case *rsa.PublicKey:
		rsaKey = k
	default:
		return ErrInvalidKey
	}


	// Create hasher
	if !m.Hash.Available() {
		return ErrHashUnavailable
	}
	hasher := m.Hash.New()
	hasher.Write([]byte(signingString))

	// Verify the signature
	err = rsa.VerifyPKCS1v15(rsaKey, m.Hash, hasher.Sum(nil), sig)
    fmt.Printf("%#v\n, %#v\n, %#v\n, %#v\n, %#v\n", rsaKey, m.Hash, hasher.Sum(nil), sig, err)
	return rsa.VerifyPKCS1v15(rsaKey, m.Hash, hasher.Sum(nil), sig)
}

// Implements the Sign method from SigningMethod
// For this signing method, must be either a PEM encoded PKCS1 or PKCS8 RSA private key as
// []byte, or an rsa.PrivateKey structure.
func (m *SigningMethodRSA) Sign(signingString string, key interface{}) (string, error) {
	var err error
	var rsaKey *rsa.PrivateKey

	switch k := key.(type) {
	case []byte:
		if rsaKey, err = ParseRSAPrivateKeyFromPEM(k); err != nil {
			return "", err
		}
	case *rsa.PrivateKey:
		rsaKey = k
	default:
		return "", ErrInvalidKey
	}

	// Create the hasher
	if !m.Hash.Available() {
		return "", ErrHashUnavailable
	}

	hasher := m.Hash.New()
	hasher.Write([]byte(signingString))

	// Sign the string and return the encoded bytes
	if sigBytes, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, m.Hash, hasher.Sum(nil)); err == nil {
		return EncodeSegment(sigBytes), nil
	} else {
		return "", err
	}
}

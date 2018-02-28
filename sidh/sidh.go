// Package sidh contains all the interfaces for the implementation of an instance 
// of Supersingular Isogeny Diffie-Hellman key exchange. The provided APIs are generic
// and they can be used by any specific implementation of SIDH over a specific prime
// for a defined quantum security level.
// Author: Amir Jalali		ajalali2016@fau.edu
package sidh

import (
	"io"
)

// Bob's public key interface.
type PublicKeyBob interface{}

// Alice's public key interface.
type PublicKeyAlice interface{}

// Bob's secret key interface.
type SecretKeyBob interface{}

// Alice's secret key interface.
type SecretKeyAlice interface{}

// The main interface for SIDH key exchange.
type SIDH interface {

	// Alice's key generation function generates Alice's public-key and secret-key.
	GenerateAliceKeypair(io.Reader) (PublicKeyAlice, SecretKeyAlice, error)

	// Bob's key generation function generates Bob's public-key and secret-key.
	GenerateBobKeypair(io.Reader) (PublicKeyBob, SecretKeyBob, error)

	// Alice's public key function generates Alice's public-key corresponds to her secret-key.
	PublicKeyAlice(SecretKeyAlice) PublicKeyAlice

	// Bob's public key function generates Bob's public-key corresponds to his secret-key
	PublicKeyBob(SecretKeyBob) PublicKeyBob

	// Alice's shared secret function generates the shared-secret using Alice's secret-key and Bob's public-key
	SharedSecretAlice(SecretKeyAlice, PublicKeyBob) []byte

	// Bob's shared secret function generates the shared-secret using Bob's secret-key and Alice's public-key
	SharedSecretBob(SecretKeyBob, PublicKeyAlice) []byte
}

// Package sike contains all the required interfaces for running an instance of
// Supersingular Isogeny Key Encapsulation (SIKE) mechanism between two parties.
// This package provide all the required APIs and generic data-types for SIKE 
// implementation without any specific underlying parameter.
// Author: Amir Jalali		ajalali2016@fau.edu
package sike

import (
	"io"
)

// SecretKey type for the key-encapsulation mechanism.
type SecretKey interface {}

// PublicKey type for the key-encapsulation mechanism.
type PublicKey interface {}

// Ciphertext type for the key-encapsulation mechanism.
type Ciphertext interface {}

// SIKE interface contains all the required APIs for running an instance of SIKE.
type SIKE interface {

	// GenerateKeyPair generates public-key and secret-key for the key-encapsulation mechanism.
	GenerateKeyPair(io.Reader) (PublicKey, SecretKey, error)

	// Encapsulation generates the ciphertext and a shared-secret.
	// The generated ciphertext is sent to the other party.
	Encapsulation(io.Reader, PublicKey) (Ciphertext, []byte, error)

	// Decapsulation generates the shared-secret using the receiver's secret-key.
	Decapsulation(SecretKey, Ciphertext) []byte
}

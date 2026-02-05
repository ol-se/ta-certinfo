// Package internal contains entities that are used through the application
package internal

import (
	"errors"
	"time"
)

var (
	// ErrParsingCert is an error returned in case of unsuccessful parsing.
	ErrParsingCert = errors.New("error parsing certificate")
	// ErrPullingCert is an error returned in case of unsuccessful pulling.
	ErrPullingCert = errors.New("error pulling certificate")
)

type (
	// CertData contains information about the certificate.
	CertData struct {
		Sub string
		Iss string
		Eat time.Time
	}

	// IDs is a helper to pass input data through the application.
	IDs struct {
		DaID string
		CID  string
	}
)

// CertDataSliceEqual returns true if two CertData slices are equal.
func CertDataSliceEqual(f, s []CertData) bool {
	if len(f) != len(s) {
		return false
	}

	for i, v := range f {
		if !CertDataEqual(v, s[i]) {
			return false
		}
	}

	return true
}

// CertDataEqual returns true if two CertData objects are equal.
func CertDataEqual(f, s CertData) bool {
	if f.Iss != s.Iss || f.Sub != s.Sub || !f.Eat.Equal(s.Eat) {
		return false
	}

	return true
}

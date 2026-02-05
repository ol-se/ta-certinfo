// Package x509parser contains x509 parsing functionality.
package x509parser

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/ol-se/ta-certinfo/internal"
)

const blockTypeCertificate = "CERTIFICATE"

// Parser is a method receiver for the x509 certificate parser.
type Parser struct{}

// Parse parses X509 certificates from *.pem. It returns parsed data or an error.
func (p *Parser) Parse(data []byte) ([]internal.CertData, error) {
	var certData []internal.CertData

	for pemBlock, rest := pem.Decode(data); pemBlock != nil; pemBlock, rest = pem.Decode(rest) {
		if pemBlock.Type != blockTypeCertificate {
			continue
		}

		cert, err := x509.ParseCertificate(pemBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", internal.ErrParsingCert, err)
		}

		certData = append(certData, internal.CertData{
			Sub: cert.Subject.String(),
			Iss: cert.Issuer.String(),
			Eat: cert.NotAfter,
		})
	}

	return certData, nil
}

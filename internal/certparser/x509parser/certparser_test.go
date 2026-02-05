package x509parser

import (
	_ "embed"
	"errors"
	"testing"
	"time"

	"github.com/ol-se/ta-certinfo/internal"
)

var (
	//go:embed testassets/good.pem
	certGood []byte

	//go:embed testassets/bad.pem
	certBad []byte
)

func TestParse(t *testing.T) {
	t.Parallel()

	t.Run("Parse: OK", func(t *testing.T) {
		t.Parallel()

		parser := Parser{}
		expectedData := []internal.CertData{
			{
				Sub: "CN=Leaf Cert 1",
				Iss: "CN=Test CA",
				Eat: time.Unix(33306230834, 0),
			},
			{
				Sub: "CN=Leaf Cert 2",
				Iss: "CN=Test CA",
				Eat: time.Unix(33306230896, 0),
			},
		}

		certData, err := parser.Parse(certGood)
		if err != nil {
			t.Errorf("Expected no error, got %v\n", err)
		}

		if !internal.CertDataSliceEqual(expectedData, certData) {
			t.Errorf("Expected %v, got %v\n", certData, expectedData)
		}
	})
	t.Run("Parse: failed to parse", func(t *testing.T) {
		t.Parallel()

		parser := Parser{}

		certData, err := parser.Parse(certBad)

		if !errors.Is(err, internal.ErrParsingCert) {
			t.Errorf("Expected error: %v, got %v\n", internal.ErrParsingCert, err)
		}

		if certData != nil {
			t.Errorf("Expected no data, got %v\n", certData)
		}
	})
}

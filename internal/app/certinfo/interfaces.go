package certinfo

import "github.com/ol-se/ta-certinfo/internal"

type (
	// Storage represents an abstract certificate storage.
	Storage interface {
		PullCert(e internal.IDs) ([]byte, error)
	}

	// Parser represents an abstract certificate parser.
	Parser interface {
		Parse(data []byte) ([]internal.CertData, error)
	}
)

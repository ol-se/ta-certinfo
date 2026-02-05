package httppuller

import (
	"io"
	"net/http"
)

// Getter represents an abstract http getter.
type Getter interface {
	Get(url string) (resp *http.Response, err error)
}

// ReadCloser is an alias type used for mock generation.
type ReadCloser = io.ReadCloser

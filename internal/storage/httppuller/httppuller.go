// Package httppuller contains functionality for pulling certificates from http storage.
package httppuller

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ol-se/ta-certinfo/internal"
)

// Puller is a method receiver for the http certificate puller.
type Puller struct {
	Hostname  string
	ReadLimit int64
	Getter    Getter
}

const (
	partDaID = "daid"
	partCID  = "cid"
)

// PullCert pulls certificates from http storage. It returns pulled data or an error.
func (p *Puller) PullCert(e internal.IDs) ([]byte, error) {
	link, err := p.buildLink(e)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", internal.ErrPullingCert, err)
	}

	resp, err := p.Getter.Get(link)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", internal.ErrPullingCert, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: invalid status code: %d", internal.ErrPullingCert, resp.StatusCode)
	}

	limitedReader := &io.LimitedReader{
		R: resp.Body,
		N: p.ReadLimit,
	}

	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", internal.ErrPullingCert, err)
	}

	return body, nil
}

func (p *Puller) buildLink(e internal.IDs) (string, error) {
	return url.JoinPath(p.Hostname, partDaID, e.DaID, partCID, e.CID)
}

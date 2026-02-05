// Package certinfo contains core logic of the application.
package certinfo

import "github.com/ol-se/ta-certinfo/internal"

// App is a method receiver for the core logic of the application.
type App struct {
	Storage Storage
	Parser  Parser
}

// PullAndParse pulls certificates from the storage. It returns parsed data or an error.
func (a *App) PullAndParse(e internal.IDs) ([]internal.CertData, error) {
	cert, err := a.Storage.PullCert(e)
	if err != nil {
		return nil, err
	}

	return a.Parser.Parse(cert)
}

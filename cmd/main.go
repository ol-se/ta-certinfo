// Package main is the entry point of the application.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ol-se/ta-certinfo/internal"
	"github.com/ol-se/ta-certinfo/internal/app/certinfo"
	"github.com/ol-se/ta-certinfo/internal/certparser/x509parser"
	"github.com/ol-se/ta-certinfo/internal/storage/httppuller"
)

// parseValidateFlags parses and validates flags.
func parseValidateFlags() (parameters, bool) {
	fHostname := flag.String("h", hostname, "hostname")
	fReadLimit := flag.Int("limit", readLimit, "read limit (bytes)")
	fDAID := flag.String("daid", "", "DAID")
	fCID := flag.String("cid", "", "CID")
	fOpt := flag.String("o", "", "iss, eat or sub (optional)")
	fIdx := flag.Int("i", -1, "index of certificate (optional, non-negative)")

	flag.Parse()

	o := opt(*fOpt)

	if !o.valid() {
		fmt.Println("Allowed -o values: iss, eat or sub")

		return parameters{}, false
	}

	if *fReadLimit < 0 {
		fmt.Println("Only positive read limit allowed")

		return parameters{}, false
	}

	if *fDAID == "" {
		fmt.Println("Missing flag: -daid")

		return parameters{}, false
	}

	if *fCID == "" {
		fmt.Println("Missing flag: -cid")

		return parameters{}, false
	}

	return parameters{
		daID:      *fDAID,
		cID:       *fCID,
		opt:       o,
		idx:       *fIdx,
		readLimit: *fReadLimit,
		hostname:  *fHostname,
	}, true
}

func printCertData(d internal.CertData) {
	fmt.Printf("Issuer: %s\nExpires at: %s\nSubject: %s\n", d.Iss, d.Eat.String(), d.Sub)
}

func printCertDataWithOpt(d internal.CertData, o opt) {
	switch o {
	case optIss:
		fmt.Println(d.Iss)
	case optEat:
		fmt.Println(d.Eat.String())
	case optSub:
		fmt.Println(d.Sub)
	}
}

func main() {
	p, ok := parseValidateFlags()
	if !ok {
		os.Exit(codeErr)
	}

	app := certinfo.App{
		Storage: &httppuller.Puller{
			Hostname:  p.hostname,
			ReadLimit: int64(p.readLimit),
			Getter:    http.DefaultClient,
		},
		Parser: &x509parser.Parser{},
	}

	certData, err := app.PullAndParse(internal.IDs{
		DaID: p.daID,
		CID:  p.cID,
	})
	if err != nil {
		log.Println(err)

		os.Exit(codeErr)
	}

	display(p, certData)

	os.Exit(codeOK)
}

func display(p parameters, certData []internal.CertData) {
	if len(certData) == 0 {
		return
	}

	switch {
	case p.idx >= len(certData):
		log.Printf("Index %d is too big (valid range: [0-%d])", p.idx, len(certData)-1)

		os.Exit(codeErr)
	case p.opt != "" && p.idx >= 0:
		printCertDataWithOpt(certData[p.idx], p.opt)
	case p.opt != "":
		for _, v := range certData {
			printCertDataWithOpt(v, p.opt)
		}
	case p.idx >= 0:
		printCertData(certData[p.idx])
	default:
		fmt.Printf("CID: %s\n", p.cID)
		fmt.Printf("DAID: %s\n", p.daID)

		for _, v := range certData {
			fmt.Println()

			printCertData(v)
		}
	}
}

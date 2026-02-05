# Certinfo
Certinfo is a utility that can parse certificates from a dedicated storage. Currently it is capable to fetch X.509 certificates from the HTTP server and display important information about them[^1].

## Go version
1.25+

## Build (Linux)
```
make build
```

## Usage
As a *.pem file can contain multiple certificates, different modes are provided for convenience.
### Default mode
Flags `-daid` and `-cid` are mandatory, no other flags are needed to use the utility in the default mode:
```bash
certinfo -daid QCDEMO -cid 3
```
Example output:
```
CID: 3
DAID: QCDEMO

Issued: CN=Test CA
Expires at: 3025-06-07 18:47:14 +0000 UTC
Subject: CN=Leaf Cert 1

Issued: CN=Test CA
Expires at: 3025-06-07 18:48:16 +0000 UTC
Subject: CN=Leaf Cert 2
```
### Single certificate
A certificated index (zero based) can be specified with the `-i` flag:
```bash
certinfo -daid QCDEMO -cid 3 -i 0
```
Example output:
```
Issued: CN=Test CA
Expires at: 3025-06-07 18:47:14 +0000 UTC
Subject: CN=Leaf Cert 1
```
### Specific fields
A field can be specified with the `-o` flag (supported values are `iss`, `sub` and `eat`):
```bash
certinfo -daid QCDEMO -cid 3 -i 0 -o sub
```
Example output:
```
CN=Leaf Cert 1
```

## Testing
Firstly, the mocks have to be generated with [mockery](https://github.com/vektra/mockery). Then:
```bash
make mock
make test
```
It is also possible to check unit tests coverage:
```bash
make cover
```

## Lint
[golangci-lint](https://github.com/golangci/golangci-lint) must be installed first.
```bash
make lint
```

[^1]: Assuming that each certificate has cID and daID, and that the important data includes Subject, Issuer and Expiration time.
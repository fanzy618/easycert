# easycert

`easycert` is a lightweight command-line utility for generating and inspecting X.509 certificates. It aims to simplify common certificate management tasks such as creating a certificate authority (CA), issuing certificates, and examining existing certificates.

## Installation

```bash
go install github.com/fanzy618/easycert@latest
```

The `easycert` binary will be installed to your `GOBIN` (usually `$GOPATH/bin`).

## Global Flags

The following flags are available to all commands:

- `-n, --name` – base name of files to read or write (default: `ca`)
- `-d, --dir` – directory where certificates and keys are stored (default: current directory)
- `--cn` – common name for the certificate (default: `easycert`)
- `--orgs` – comma-separated list of organizations (default: `easycert`)
- `--dns` – repeatable flag for DNS names
- `--ip` – repeatable flag for IP addresses
- `-b` – RSA key size in bits (default: `2048`)

## Commands

### `ca`
Generate a self-signed certificate authority.

```bash
easycert ca -n root --cn root-ca
```

This command writes `root.pem` and `root-key.pem` to the directory specified by `--dir`.

### `cert`
Generate a certificate, optionally signed by an existing CA.

```bash
# create a certificate signed by the previously generated root
 easycert cert -n cert --cn mysite.com --ca root --dns mysite1.com --dns mysite2.com
```

If `--ca` is omitted a new self-signed CA is generated automatically.

### `show`
Display information about a certificate stored on disk.

```bash
easycert show -n cert
```

Example output:

```
Subject: CN=mysite.com,O=easycert
DNSNames: mysite1.com, mysite2.com
IPAddresses:
NotBefore: 2025-08-07T06:46:02Z
NotAfter: 2035-08-05T06:46:04Z
```

## License

This project is licensed under the Apache 2.0 License.


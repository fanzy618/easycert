/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cert

import (
	"bytes"
	"crypto"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"path"
	"time"

	certutil "k8s.io/client-go/util/cert"
	"k8s.io/client-go/util/keyutil"
)

const (
	KeyExt  = ".key"
	CertExt = ".crt"
)

type Config struct {
	certutil.Config

	Name string
	Dir  string

	Bits     int
	ValidFor time.Duration
}

// NewKey return a new private key
func NewKey(algorithm string, bits int) (crypto.Signer, error) {
	key, err := rsa.GenerateKey(cryptorand.Reader, bits)
	if err != nil {
		log.Println("NewKey failed:", err.Error())
		return nil, err
	}
	return key, nil
}

func LoadFromDisk(pkiPath, name string) (*x509.Certificate, crypto.Signer, error) {
	certPath := path.Join(pkiPath, name+CertExt)
	certs, err := certutil.CertsFromFile(certPath)
	if err != nil {
		log.Println("LoadFromDisk: Load cert fialed.", pkiPath, name, err)
		return nil, nil, err
	}

	keyPath := path.Join(pkiPath, name+KeyExt)
	key, err := keyutil.PrivateKeyFromFile(keyPath)
	if err != nil {
		log.Println("LoadFromDisk: Load key failed.", pkiPath, name, err)
		return nil, nil, err
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		log.Println("LoadFromDisk: Load key failed, Unknown key format")
		return nil, nil, fmt.Errorf("Unknown key format")
	}
	return certs[0], rsaKey, nil
}

// NewSelfSignedCACert creates a self-signed CA certificate respecting the
// validity duration configured in cfg. If cfg.ValidFor is zero, a default of
// one year is used.
func NewSelfSignedCACert(cfg *Config, key crypto.Signer) (*x509.Certificate, error) {
	now := time.Now()
	validFor := cfg.ValidFor
	if validFor == 0 {
		validFor = time.Hour * 24 * 365
	}
	tmpl := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		NotBefore:             now.UTC(),
		NotAfter:              now.Add(validFor).UTC(),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	certDERBytes, err := x509.CreateCertificate(cryptorand.Reader, &tmpl, &tmpl, key.Public(), key)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}

func WriteToDisk(pkiPath, name string, certs []*x509.Certificate, key crypto.Signer) error {

	keyInPem, err := keyutil.MarshalPrivateKeyToPEM(key)
	if err != nil {
		log.Println("MarshalPrivateKeyToPEM", err)
		return err
	}
	keyPath := path.Join(pkiPath, name+KeyExt)
	err = keyutil.WriteKey(keyPath, keyInPem)
	if err != nil {
		log.Println("WriteKey", err)
		return err
	}
	log.Println("Write key to:\t", keyPath)

	certPath := path.Join(pkiPath, name+CertExt)
	pemData := &bytes.Buffer{}
	for _, cert := range certs {
		err := pem.Encode(pemData, &pem.Block{Type: certutil.CertificateBlockType, Bytes: cert.Raw})
		if err != nil {
			log.Println("PEM Encode error:", err)
			return err
		}
	}

	err = certutil.WriteCert(certPath, pemData.Bytes())
	if err != nil {
		log.Println("WriteCert", err)
		return err
	}
	log.Println("Write cert to:\t", certPath)
	return nil
}

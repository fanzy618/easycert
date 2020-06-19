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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"path"

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

	Bits int
}

// NewKey return a new private key
func NewKey(algorithm string, bits int) (crypto.Signer, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Println("NewKey failed:", err.Error)
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

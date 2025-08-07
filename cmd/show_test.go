package cmd

import (
	"bytes"
	"crypto/x509"
	"net"
	"strings"
	"testing"

	certutil "k8s.io/client-go/util/cert"

	"github.com/fanzy618/easycert/pkg/cert"
)

func TestShowCmd(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &cert.Config{
		Config: certutil.Config{
			CommonName: "test",
			AltNames: certutil.AltNames{
				DNSNames: []string{"example.com"},
				IPs:      []net.IP{net.ParseIP("127.0.0.1")},
			},
		},
		Name: "test",
		Dir:  tmpDir,
		Bits: 2048,
	}

	key, err := cert.NewKey("", cfg.Bits)
	if err != nil {
		t.Fatalf("NewKey: %v", err)
	}
	caCert, err := certutil.NewSelfSignedCACert(cfg.Config, key)
	if err != nil {
		t.Fatalf("NewSelfSignedCACert: %v", err)
	}
	if err := cert.WriteToDisk(cfg.Dir, cfg.Name, []*x509.Certificate{caCert}, key); err != nil {
		t.Fatalf("WriteToDisk: %v", err)
	}

	globalCfg = cfg
	buf := new(bytes.Buffer)
	showCmd.SetOut(buf)
	if err := showCmd.RunE(showCmd, []string{}); err != nil {
		t.Fatalf("RunE: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Subject:") {
		t.Fatalf("expected output to contain Subject, got %s", out)
	}
}

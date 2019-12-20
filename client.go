package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"
)

func verifyPeerCertificate(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	if len(rawCerts) != 1 {
		return fmt.Errorf("Expected 1 certificate, but got %d", len(rawCerts))
	}

	cert, err := x509.ParseCertificate(rawCerts[0])

	if err != nil {
		return err
	}

	currentTime := time.Now().UTC()

	if currentTime.Before(cert.NotBefore) {
		return fmt.Errorf("Cert is not valid before %v and time is %v", cert.NotBefore, currentTime)
	}

	if currentTime.After(cert.NotAfter) {
		return fmt.Errorf("Cert is not valid after %v and time is %v", cert.NotAfter, currentTime)
	}

	return cert.VerifyHostname("localhost")
}

func main() {
	conn, err := tls.Dial("tcp", "127.0.0.1:443", &tls.Config{
		InsecureSkipVerify:    true,
		VerifyPeerCertificate: verifyPeerCertificate,
	})

	if err != nil {
		panic("failed to connect: " + err.Error())
	}

	conn.Close()
}

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"time"
)

type OutOfBandHandler struct {
	Certificate string
}

func (oob OutOfBandHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, oob.Certificate)
}

func main() {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	notBefore := time.Now().UTC().Add(-(time.Minute * 5))
	notAfter := notBefore.Add(time.Hour * 5)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %s", err)
	}

	template := x509.Certificate{
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:     []string{"localhost"},
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Some Enterprise"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDer, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}

	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDer})

	//Simulate Out Of Band Handler (obviously this makes no sense to be on the same transport, this
	//is just for demonstration purposes)
	go http.ListenAndServe("localhost:8080", &OutOfBandHandler{
		Certificate: string(certPem),
	})

	privDer, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatalf("Unable to marshal private key: %v", err)
	}

	privPem := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDer})

	cert, err := tls.X509KeyPair(certPem, privPem)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	srv := &http.Server{
		Addr:         "localhost:443",
		TLSConfig:    cfg,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	//Create TLS transport
	log.Fatal(srv.ListenAndServeTLS("", ""))
}

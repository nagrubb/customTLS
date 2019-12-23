package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getOobCert() []byte {
	rsp, err := http.Get("http://localhost:8080")

	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()

	certPem, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		panic(err)
	}

	return certPem
}

func main() {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		panic(err)
	}

	certPool.AppendCertsFromPEM(getOobCert())

	conn, err := tls.Dial("tcp", "localhost:443", &tls.Config{
		RootCAs: certPool,
	})

	if err != nil {
		panic("failed to connect: " + err.Error())
	}

	conn.Close()
	fmt.Println("Success!")
}

/*
 * Copyright 2022 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// A fake & minimal https service that establishes a short-lived connection
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"time"
)

func main() {

	var certificate = GenerateSSLCert()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Currently visiting home page of a HTTPS based service")
	})

	log.Println("Starting server at https://127.0.0.1:8443/")
	var err = ListenAndServeTLSKeyPair(":8443", certificate, mux)
	if err != nil {
		log.Fatalln(err)
	}

}

// GenerateSSLCert is responsible for providing a self signed certificate
func GenerateSSLCert() tls.Certificate {

	// Generating a private key which will also be used to sign the certificate
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// Each certificate must have a unique serial number.
	// Here a randomly generated 128 bit number is used.
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	// Specifying the template for the self signed certificate.
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Tsunami Testbed"},
		},
		DNSNames:  []string{"localhost"},
		NotBefore: time.Now(),
		// Certificate is valid for 10 years after creation
		NotAfter: time.Now().Add(10 * 365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Creating a certificate using the template specified.
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	var SignedCert tls.Certificate
	SignedCert.Certificate = append(SignedCert.Certificate, certBytes)
	SignedCert.PrivateKey = privateKey
	return SignedCert
}

// ListenAndServeTLSKeyPair starts a server using the TLS certificate.
// Using our own implementation for ListenAndServeTLSKeyPair function because
// the library version takes file path of the certificate and private key as parameters
// whereas we are using an in-memory certificate.
func ListenAndServeTLSKeyPair(address string, cert tls.Certificate, handler http.Handler) error {

	if address == "" {
		return errors.New("Invalid address string")
	}
	server := &http.Server{Addr: address, Handler: handler}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"http/1.1"},
		MinVersion:   tls.VersionTLS13,
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(tcpKeepAliveListener{listener.(*net.TCPListener)},
		config)

	return server.Serve(tlsListener)

}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (listener tcpKeepAliveListener) Accept() (con net.Conn, err error) {
	tc, err := listener.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(false)
	return tc, nil
}

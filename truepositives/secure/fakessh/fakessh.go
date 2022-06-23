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

// A fake & minimal ssh service that establishes a short-lived connection when the specified
// user and password is used.
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

var (
	sshUser     = flag.String("user", "root", "Allowed user for the fake ssh service.")
	sshPassword = flag.String("password", "qwerty", "Configured password for the fake ssh service.")
	sshPort     = flag.Int("port", 8022, "Default port to run ssh server at.")
)

func main() {
	flag.Parse()

	config := &ssh.ServerConfig{
		// Check if user & password matches with the pre-defined values
		PasswordCallback: func(c ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			if c.User() == *sshUser && string(password) == *sshPassword {
				return nil, nil
			}
			return nil, fmt.Errorf("Failed authentication for user %q password %q from %q", c.User(), string(password), c.RemoteAddr().String())
		},
	}

	config.AddHostKey(generateHostKey())

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *sshPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %s", *sshPort, err)
	}

	log.Printf("Listening on port %d", *sshPort)

	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming tcp connection: %s", err)
			continue
		}

		sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, config)
		if err != nil {
			log.Printf("Failed to establish initial handshake: %s", err)
			continue
		}

		log.Printf("Established a new SSH connection from %s %s", sshConn.RemoteAddr(), sshConn.ClientVersion())
		// Discard all global out-of-band Requests
		go ssh.DiscardRequests(reqs)
		// Accept all channels
		go handleChannels(chans)
	}
}

func handleChannels(chans <-chan ssh.NewChannel) {
	for newChannel := range chans {
		go handleChannel(newChannel)
	}
}

func handleChannel(newChannel ssh.NewChannel) {
	connection, requests, err := newChannel.Accept()
	if err != nil {
		log.Printf("Failed to accept channel %s", err)
		return
	}

	// Reply & close the connection right away
	go func() {
		for req := range requests {
			log.Printf("[%s] type request received with payload [%q]", req.Type, string(req.Payload))
			if len(req.Payload) == 0 {
				req.Reply(true, nil)
			}
			connection.Close()
			log.Printf("Session closed")
		}
	}()

}

func generateHostKey() ssh.Signer {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatalf("Failed to generate key: %s", err)
	}
	privateKeyPkcs := x509.MarshalPKCS1PrivateKey(privateKey)
	privateBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyPkcs,
	}
	privatePEM := pem.EncodeToMemory(&privateBlock)
	signer, err := ssh.ParsePrivateKey(privatePEM)
	if err != nil {
		log.Fatalf("Failed to parse private key: %s", err)
	}
	return signer
}

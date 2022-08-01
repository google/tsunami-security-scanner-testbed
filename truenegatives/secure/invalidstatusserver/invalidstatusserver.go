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

// A non-standard http server that returns invalid http status.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

var (
	serverPort = flag.Int("port", 8080, "Default port to run server at.")
)

func main() {
	server, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *serverPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %s", *serverPort, err)
	}

	defer server.Close()
	log.Printf("Listening on port %d", *serverPort)

	for {
		tcpConn, err := server.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming tcp connection: %s", err)
			continue
		}
		connAddr := tcpConn.RemoteAddr()
		if connAddr.String() == "" {
			connAddr = tcpConn.LocalAddr()
		}
		log.Printf("Established a new tcp connection from %s", connAddr)
		go processRequest(tcpConn)
	}
}

func processRequest(connection net.Conn) {
	// skip reading the request, just reply with an invalid HTTP status
	connection.Write([]byte("HTTP/1.0 001 \n\nHello"))
	connection.Close()
}

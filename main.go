package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
)

var (
	listenPort        = flag.Int("l", 443, "The port to listen")
	certificationFile = flag.String("c", "", "Certification file")
	privateKeyFile    = flag.String("k", "", "Private Key file")
	verbose           = flag.Bool("v", false, "Verbose Output")
)

func main() {
	flag.Parse()
	var (
		s        net.Listener
		certPair tls.Certificate
		err      error
	)

	if *certificationFile != "" && *privateKeyFile != "" {
		certPair, err = tls.LoadX509KeyPair(*certificationFile, *privateKeyFile)
		if err != nil {
			log.Fatalln("Failed to parse certificate", err)
		}
	}

	config := tls.Config{Certificates: []tls.Certificate{certPair}}
	s, err = tls.Listen("tcp", fmt.Sprintf(":%d", *listenPort), &config)
	if err != nil {
		panic(err)
	}

	log.Println("Started VPN/Sentry server..")

	con, _ := s.Accept()
	defer con.Close()

	// Read initial client packet and write it back exactly the same
	b := make([]byte, 1500)
	n, err := con.Read(b)
	if err != nil {
		log.Fatalln("Client read failed", err)
	}
	log.Println("Client read:", string(b))
	log.Println("Server send:", string(b))
	_, err = con.Write(b[:n])
	if err != nil {
		log.Fatalln("Server send failed", err)
	}

	// Send hardcoded message back to client
	val := []byte{'\x17', '\x00', '\x02', '\x05', '\xdc'}
	log.Println("Server send:", val)
	_, err = con.Write(val)
	if err != nil {
		log.Fatalln("Server send failed", err)
	}

	// Send hardcoded message back to client
	val = []byte{'\x16', '\x00', '\x0c', '\x00', '\x05', '\xfd', '\xf7', '\x00', '\x01', '\x18', '\x88', '\x06', '\x40', '\x00', '\x00'}
	log.Println("Server send:", val)
	_, err = con.Write(val)
	if err != nil {
		log.Fatalln("Server send failed", err)
	}

	// Read 5 responses from client (about profiles, etc)
	for i := 0; i < 5; i++ {
		b = make([]byte, 1500)
		n, err = con.Read(b)
		if n == 0 || err != nil {
			log.Fatalln("Client read failed", err)
		}
		log.Println("Client read:", string(b))
		log.Println("Client read (bytes):", b[:n])
	}

	// Send hardcoded message back to client
	val = []byte{'\x02', '\x00', '\x05', '\x04', '\x6d', '\x61', '\x69', '\x6e'}
	log.Println("Server send:", val)
	_, err = con.Write(val)
	if err != nil {
		log.Fatalln("Server send failed", err)
	}

	// Send hardcoded message back to client
	val = []byte{'\x16', '\x00', '\x0c', '\x00', '\x05', '\xfd', '\xf7', '\x00', '\x01', '\x18', '\x88', '\x06', '\x40', '\x00', '\x00'}
	log.Println("Server send:", val)
	_, err = con.Write(val)
	if err != nil {
		log.Fatalln("Server send failed", err)
	}

	// Should be connected now - need to figure out how to proxy data
	// For now just read until connection closes
	for {
		b = make([]byte, 1500)
		n, err := con.Read(b)
		if n == 0 || err != nil {
			log.Fatal("Connection closed")
		}
		log.Println("Client read:", string(b))
		log.Println("Client read (bytes):", b[:n])
	}
}

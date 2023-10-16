package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
)

func handleConnection(src net.Conn, destPort string) {
	dest, err := net.Dial("tcp", "localhost:"+destPort)
	if err != nil {
		log.Printf("Could not connect to localhost:%s: %v", destPort, err)
		return
	}
	defer dest.Close()

	// Extract source and destination IPs
	srcIP := src.RemoteAddr().(*net.TCPAddr).IP.String()
	destIP := dest.LocalAddr().(*net.TCPAddr).IP.String()

	// Resolve hostnames
	srcHostnames, err := net.LookupAddr(srcIP)
	if err != nil || len(srcHostnames) == 0 {
		srcHostnames = []string{srcIP}
	}

	destHostnames, err := net.LookupAddr(destIP)
	if err != nil || len(destHostnames) == 0 {
		destHostnames = []string{destIP}
	}

	go func() {
		buf := new(bytes.Buffer)
		_, err := io.Copy(buf, src)
		if err != nil {
			log.Printf("Could not copy tcp source: %v", err)
			return
		}
		payload := buf.String()

		// Create log entry
		logEntry := map[string]interface{}{
			"source_ip":            srcIP,
			"destination_ip":       destIP,
			"source_hostname":      srcHostnames[0],  // Taking the first resolved name
			"destination_hostname": destHostnames[0], // Taking the first resolved name
			"tcp_payload":          payload,
		}

		// Convert map to JSON
		data, err := json.Marshal(logEntry)
		if err != nil {
			log.Printf("Error marshalling log entry: %v", err)
			return
		}
		log.Println("POD-CONN\n", string(data)) // Send to stdout so fluentd can "catch" it

		// Forward the buffered data to the backend
		_, err = io.Copy(dest, buf)
		if err != nil {
			log.Printf("Error while forwarding buffered data to destination: %v", err)
		}
	}()

	go func() {
		_, err := io.Copy(src, dest)
		if err != nil {
			log.Printf("Error while copying data from destination to source: %v", err)
		}
	}()
}

func main() {
	listenPort := os.Getenv("LISTEN_PORT")
	if listenPort == "" {
		log.Fatalf("LISTEN_PORT environment variable not set")
	}

	listener, err := net.Listen("tcp", ":"+listenPort)
	if err != nil {
		log.Fatalf("Error while starting server: %v", err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatalf("Couldn't close network listener")
		}
	}(listener)

	destPort := os.Getenv("DEST_PORT")
	if destPort == "" {
		log.Fatalf("DEST_PORT environment variable not set")
	}

	log.Printf("Server started on :%s, forwarding to localhost:%s", listenPort, destPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error while accepting connection: %v", err)
			continue
		}
		go handleConnection(conn, destPort)
	}
}

// Package main provides a minimal TCP health check binary for scratch-based
// Docker images. It performs a TCP dial to the local service port and exits
// with code 0 on success or 1 on failure.
//
// Uses raw net.Dial instead of net/http to avoid pulling in TLS/crypto
// packages, resulting in a much smaller binary (~2MB vs ~5.6MB).
package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	addr := "localhost:8080"
	if envAddr := os.Getenv("HEALTHCHECK_ADDR"); envAddr != "" {
		addr = envAddr
	}

	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck failed: %v\n", err)
		os.Exit(1)
	}
	if err := conn.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to close connection: %v\n", err)
		os.Exit(1)
	}
}

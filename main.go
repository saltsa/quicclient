package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

var t *time.Ticker
var certPool *x509.CertPool

func init() {
	certPool = x509.NewCertPool()
}

func client(addr string) error {

	if addr == "" {
		log.Println("empty addr")
		return errors.New("empty addr")
	}
	var qconf quic.Config

	qconf.TokenStore = quic.NewLRUTokenStore(100, 200)
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs:            certPool,
			InsecureSkipVerify: true,
		},
		QuicConfig: &qconf,
	}

	hc := http.Client{
		Timeout:   15 * time.Second,
		Transport: roundTripper,
	}

	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		log.Fatalln(err)
	}

	t = time.NewTicker(5 * time.Second)

	for range t.C {
		log.Printf("clientreq: %s", addr)
		start := time.Now()
		resp, err := hc.Do(req)
		if err != nil {
			log.Printf("client error: %s", err)
			continue
		}
		log.Printf("status %d in %s", resp.StatusCode, time.Since(start))
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("body read error: %s", err)
		}
		resp.Body.Close()

		log.Printf("body: %s", body)
	}

	log.Println("client closing")
	return nil
}

func main() {
	flag.Parse()

	addr := flag.Arg(0)

	go client(addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		recv := <-quit
		if recv == os.Interrupt {
			os.Stdout.Write([]byte("\r"))
		}

		log.Printf("got signal: %v", recv)

		if recv == os.Interrupt || recv == syscall.SIGTERM {
			log.Println("quitting")
			return
		}
	}

}

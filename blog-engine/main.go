package main

import (
	"blog-engine/article"
	"log"
	"net"
	"os"
	"strconv"
)

func configureHost(s *article.HTTPServer) error {
	hostStr, ok := os.LookupEnv("ARTICLE_HOST")
	if !ok {
		return nil
	}

	if ip := net.ParseIP(hostStr); ip == nil {
		return nil
	}

	s.Host = hostStr
	return nil
}

func configurePort(s *article.HTTPServer) error {
	portStr, ok := os.LookupEnv("ARTICLE_PORT")
	if !ok {
		return nil
	}

	port, err := strconv.ParseUint(portStr, 10, 16)

	if err != nil {
		return err
	}

	s.Port = uint16(port)
	return nil
}

func main() {
	server, err := article.NewHTTPServer(configureHost, configurePort)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting HTTP server at %s port %d...", server.Host1(), server.Port1())
	server.Start()
}

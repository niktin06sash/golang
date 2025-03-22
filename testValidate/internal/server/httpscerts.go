package server

import (
	"crypto/tls"
	"encoding/base64"
	"log"
	"os"
	"testValidate/internal/erro"
)

func loadCertsFromEnv() (tls.Certificate, error) {
	crtBase64 := os.Getenv("SERVER_CERT")
	keyBase64 := os.Getenv("SERVER_KEY")

	if crtBase64 == "" || keyBase64 == "" {
		return tls.Certificate{}, erro.ErrorCertsEnv
	}

	crtBytes, err := base64.StdEncoding.DecodeString(crtBase64)
	if err != nil {
		log.Println(err)
		return tls.Certificate{}, erro.ErrorDecodeCert
	}

	keyBytes, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		log.Println(err)
		return tls.Certificate{}, erro.ErrorDecodeKey
	}
	cert, err := tls.X509KeyPair(crtBytes, keyBytes)
	if err != nil {
		log.Println(err)
		return tls.Certificate{}, erro.ErrorX509KeyPair
	}
	return cert, nil
}

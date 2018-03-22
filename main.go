package main

import (
	"net/http"
	"github.com/brokenbydefault/Server/handlers/Websocket"
	"github.com/brokenbydefault/Server/handlers/Callback"
	"github.com/brokenbydefault/Server/SQL"
	"github.com/brokenbydefault/Server/Apps/NanofyApp"
	"github.com/brokenbydefault/Server/Config"
	"crypto/tls"
	"log"
	"github.com/brokenbydefault/Server/Security"
)

func init() {
	Config.Start()
	SQL.Start()
	NanofyApp.Start()
}

func main() {
	api := http.NewServeMux()
	api.HandleFunc("/", Security.SetHeaders(Websocket.Start))
	srv := &http.Server{
		Addr:    Config.Config["API_IP"] + ":443",
		Handler: api,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}

	go func() {
		log.Panic(srv.ListenAndServeTLS(Config.Dir()+"/../.pki/cert.crt", Config.Dir()+"/../.pki/key.pem"))
	}()

	callback := http.NewServeMux()
	callback.HandleFunc("/", Security.SetHeaders(Callback.Start))
	// This port SHOULD be closed!
	// Because the communication occur in localhost there is no need for SSL.
	log.Panic(http.ListenAndServe(Config.Config["CALLBACK_IP"]+":7771", callback))
}

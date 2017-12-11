// Copyright 2017 Johan Brandhorst. All Rights Reserved.
// See LICENSE for licensing terms.

package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/pat"
	"github.com/gorilla/websocket"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/johanbrandhorst/protobuf/wsproxy"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	"github.com/anxiousmodernman/grpcweb-boilerplate/backend"
	"github.com/anxiousmodernman/grpcweb-boilerplate/frontend/bundle"
	"github.com/anxiousmodernman/grpcweb-boilerplate/proto/server"
)

// TODO pass this down to my object
var logger *logrus.Logger

func init() {
	logger = logrus.StandardLogger()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
		DisableSorting:  true,
	})
	// Should only be done from init functions
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(logger.Out, logger.Out, logger.Out))
}

func main() {

	// Proxy is our code that implements generated interface for server.
	prxy, err := backend.NewProxy("co-chair.db")
	if err != nil {
		log.Fatalf("proxy init: %v", err)
	}

	gs := grpc.NewServer()
	server.RegisterProxyServer(gs, prxy)
	wrappedServer := grpcweb.WrapServer(gs)

	clientCreds, err := credentials.NewClientTLSFromFile("./cert.pem", "")
	if err != nil {
		logger.WithError(err).Fatal("Failed to get local server client credentials, did you run `make generate_cert`?")
	}

	wsproxy := wsproxy.WrapServer(
		http.HandlerFunc(wrappedServer.ServeHTTP),
		wsproxy.WithLogger(logger),
		wsproxy.WithTransportCredentials(clientCreds))

	p := pat.New()
	_ = p

	handler := func(resp http.ResponseWriter, req *http.Request) {
		// Redirect gRPC and gRPC-Web requests to the gRPC-Web Websocket Proxy server
		if req.ProtoMajor == 2 && strings.Contains(req.Header.Get("Content-Type"), "application/grpc") ||
			websocket.IsWebSocketUpgrade(req) {
			wsproxy.ServeHTTP(resp, req)
		} else {
			// Serve the GopherJS client
			http.FileServer(bundle.Assets).ServeHTTP(resp, req)
		}
	}

	// auth feature flag
	var withauth = false
	var h http.Handler
	if withauth {
		jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("COCHAIR_AUTH0_SECRET")), nil
			},
			// When set, the middleware verifies that tokens are signed with the
			// specific signing algorithm If the signing method is not constant
			// the ValidationKeyGetter callback can be used to implement additional checks
			// Important to avoid security issues described here:
			// https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
			SigningMethod: jwt.SigningMethodHS256,
		})
		h = jwtMiddleware.Handler(http.HandlerFunc(handler))
	} else {
		h = http.HandlerFunc(handler)
		// https://dsasf.auth0.com/login?client=xxx
	}

	addr := "localhost:2016"
	httpsSrv := &http.Server{
		Addr:    addr,
		Handler: h,
		// Some security settings
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
		TLSConfig: &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},
		},
	}

	logger.Info("Serving on https://" + addr)
	logger.Fatal(httpsSrv.ListenAndServeTLS("./cert.pem", "./key.pem"))
}

package main

import (
	"strings"

	"honnef.co/go/js/dom"

	"github.com/anxiousmodernman/grpcweb-boilerplate/proto/client"
)

// Build this snippet with GopherJS, minimize the output and
// write it to html/frontend.js
//go:generate gopherjs build frontend.go -m -o html/frontend.js

// Integrate generated JS into a Go file for static loading.
//go:generate bash -c "go run assets_generate.go"

// This constant is very useful for interacting with
// the DOM dynamically
var document = dom.GetWindow().Document().(dom.HTMLDocument)

// Define no-op main since it doesn't run when we want it to
func main() {}

// Ensure our setup() gets called as soon as the DOM has loaded
func init() {
	document.AddEventListener("DOMContentLoaded", false, func(_ dom.Event) {
		go setup()
	})
}

// Setup is where we do the real work.
func setup() {
	// This is the address to the server, and should be used
	// when creating clients.
	serverAddr := strings.TrimSuffix(document.BaseURI(), "/")

	// TODO: Use functions exposed by generated interface
	_ = client.NewProxyClient(serverAddr)

	document.Body().SetInnerHTML(`<div><h2>GopherJS gRPC-Web is great!</h2></div>`)
}

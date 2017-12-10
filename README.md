# GopherJS gRPC-Web Boilerplate

Forked from the amazing gift Johan Brandhorst gave to the world.

A minimal repo containing all the boilerplate for getting started with
GopherJS using gRPC-Web. Clone and use as a base for your own
experimentation with GopherJS and gRPC-Web!

It comes complete with a Go generate template for statically serving
your generated JS code from the backend server.

## Usage

After cloning the repo, there are a couple of initial steps;

1. Install the generate dependencies with `make install`.
1. Generate a self-signed certificate with `make generate_cert`.
1. Edit the `go_package` and `gopherjs.gopherjs_package` protobuf options in `proto/web.proto`. These should be updated to point to the packages that are used in your repo.
1. Edit the imports in `main.go`, `backend/backend.go` and `frontend/frontend.go` to point to the the paths used in your repo.

After that is done, you can generate the JS files using `make generate`,
and try out the web server using `make serve`.

## Making it your own

The next step is to define the interface you want to expose in
`proto/web.proto`. See https://developers.google.com/protocol-buffers/
tutorials and guides on writing protofiles.

Once that is done, regenerate the backend and frontend files using
`make generate`. This will mean you'll need to implement any functions in
`backend/backend.go`, or else the build will fail since your struct won't
be implementing the interface defined by the generated file in `proto/server/`.

It also means you can start using the functions exposed by the server
in your frontend client `frontend/frontend.go`. This should hopefully be
all you need to get started playing around with the GopherJS gRPC-Web
bindings!

## What this repo isn't

This repo is _not_ a general example of how to use the GopherJS gRPC-Web bindings.
For such an example, please see https://github.com/johanbrandhorst/grpcweb-example
and https://grpcweb.jbrandhorst.com.

This repo is also not a particularly good example of how to write an app with
GopherJS. Please explore the tutorials and wiki pages on
https://github.com/gopherjs/gopherjs and check out frameworks such as
https://myitcv.io/react and https://github.com/gopherjs/vecty for
a better way to write GopherJS apps.

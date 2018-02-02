# gRPC Tutorial

From the video: [GopherCon 2017: Alan Shreve - grpc: From Tutorial to Production](https://www.youtube.com/watch?v=7FZ6ZyzGex0)



## Setup

### Install gRPC

`go get -u google.golang.org/grpc`


#### Install Protocol Buffers v3

Download pre-compiled binaries for your platform(`protoc-<version>-<platform>.zip`) from here: https://github.com/google/protobuf/releases

Update the environment variable PATH to include the path to the protoc binary file.


#### Install protoc plugin for Go

`go get -u github.com/golang/protobuf/protoc-gen-go`



### Generate gRPC code

`protoc protoc api/app.proto --go_out=plugins=grpc:.`


### Generate a self-signed X.509 certificate for a TLS server

`go run "$(go env GOROOT)/src/crypto/tls/generate_cert.go" -host localhost`


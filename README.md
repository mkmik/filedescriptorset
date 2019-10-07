# filedescriptorset
    go get github.com/mkmik/filedescriptorset

[![GoDoc](https://godoc.org/github.com/mkmik/filedescriptorset?status.png)](https://godoc.org/github.com/mkmik/filedescriptorset)

Package filedescriptorset helps you embed protobuf descriptors, and optionally their transitive dependencies.

1. add "-odeps.bin" to your protoc command. You can also add "--include_imports" if you
   want to include transitively imported descriptors.

2. create a .go file in the same dir where protoc generates code, adding this to it:

```
//go:generate go-bindata -modtime 1 -mode 420 -o deps.go -pkg client_model_proto3desc deps.bin

func init() { filedescriptorset.MustRegisterFileSet(MustAsset("deps.bin")) }
```

// Package filedescriptorset helps you embed protobuf descriptors, and optionally their transitive dependencies.
//
// 1. add "-odeps.bin" to your protoc command. You can also add "--include_imports" if you
//    want to include transitively imported descriptors.
// 2. create a .go file in the same dir where protoc generates code, adding this to it:
//
//     //go:generate go-bindata -modtime 1 -mode 420 -o deps.go -pkg client_model_proto3desc deps.bin
//
//     func init() { filedescriptorset.MustRegisterFileSet(MustAsset("deps.bin")) }
package filedescriptorset

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// RegisterFileSet takes an uncompressed binary encoded FileDescriptorSet
// (usually produced by protoc --include_imports -ofile.bin).
func RegisterFileSet(b []byte) error {
	var d descriptor.FileDescriptorSet
	if err := proto.Unmarshal(b, &d); err != nil {
		return err
	}
	for _, f := range d.File {
		fb, err := proto.Marshal(f)
		if err != nil {
			return err
		}
		var fzb bytes.Buffer
		fzw := gzip.NewWriter(&fzb)
		if _, err := fzw.Write(fb); err != nil {
			return err
		}
		fzw.Close()
		proto.RegisterFile(f.GetName(), fzb.Bytes())
	}
	return nil
}

// MustRegisterFileSet is like RegisterFileSet but panics on error.
// This can be useful in global or func init() initializers.
func MustRegisterFileSet(b []byte) {
	if err := RegisterFileSet(b); err != nil {
		panic(err)
	}
}

// RegisterCompressedFileSet is like RegisterFileSet but accepts a gzip compressed
// binary encoded FileDescriptorSet proto.
func RegisterCompressedFileSet(b []byte) error {
	r, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return err
	}
	return RegisterFileSet(buf.Bytes())
}

// MustRegisterCompressedFileSet is like RegisterCompressedFileSet but panics on error.
// This can be useful in global or func init() initializers.
func MustRegisterCompressedFileSet(b []byte) {
	if err := RegisterCompressedFileSet(b); err != nil {
		panic(err)
	}
}

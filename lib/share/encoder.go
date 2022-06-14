package share

import (
	"bytes"
	"compress/lzw"
	"encoding/base64"
	"fmt"
	"io"
)

func Encode(data string) (string, error) {
	// Compress data
	buf := &bytes.Buffer{}
	w := lzw.NewWriter(buf, lzw.LSB, 8)
	fmt.Fprint(w, data)
	if err := w.Close(); err != nil {
		return "", err
	}
	// Encode to base64
	return base64.RawURLEncoding.EncodeToString(buf.Bytes()), nil
}

func Decode(data string) (string, error) {
	// Decode base64
	b, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	// Decompress data
	r := lzw.NewReader(bytes.NewReader(b), lzw.LSB, 8)
	defer r.Close()
	res, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

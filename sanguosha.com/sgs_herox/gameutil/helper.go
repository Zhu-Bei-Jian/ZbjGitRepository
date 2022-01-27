package gameutil

import (
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strings"
)

func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	cipherStr := h.Sum(nil)
	ret := hex.EncodeToString(cipherStr)
	return ret
}
func MD5Bytes(s []byte) string {
	h := md5.New()
	h.Write(s)
	cipherStr := h.Sum(nil)
	ret := hex.EncodeToString(cipherStr)
	return ret
}
func Trim(s string) string {
	return strings.TrimSpace(s)
}

func SHA512(src string) string {
	h := sha512.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

func ZipBytes(b []byte) ([]byte, error) {
	var ret bytes.Buffer
	w := zlib.NewWriter(&ret)
	w.Write(b)
	w.Close()
	return ret.Bytes(), nil
}
func ZipBase64(b []byte) (string, error) {
	var ret bytes.Buffer
	w := zlib.NewWriter(&ret)
	w.Write(b)
	w.Close()
	return base64.StdEncoding.EncodeToString(ret.Bytes()), nil
}
func UnZipBytes(decodeBytes []byte) ([]byte, error) {
	var out bytes.Buffer
	b := bytes.NewBuffer(decodeBytes)
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes(), nil
}
func UnZipBase64(s string) ([]byte, error) {
	decodeBytes, e := base64.StdEncoding.DecodeString(s)
	if e != nil {
		return nil, e
	}
	var out bytes.Buffer
	b := bytes.NewBuffer(decodeBytes)
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes(), nil
}

package util

import (
	"io"
	"reflect"
	"unsafe"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// UTF8ToGBK UTF8 To GBK
func UTF8ToGBK(reader io.Reader) io.Reader {
	return transform.NewReader(reader, simplifiedchinese.GBK.NewEncoder())
}

// GBKToUTF8 GBK To UTF8
func GBKToUTF8(reader io.Reader) io.Reader {
	return transform.NewReader(reader, simplifiedchinese.GBK.NewDecoder())
}

// String force casts a []byte to a string.
// USE AT YOUR OWN RISK.
// (code from github.com/youtube/vitess)
func String(b []byte) (s string) {
	if len(b) == 0 {
		return ""
	}
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len
	return
}

// Slice force casts a string to a []byte.
// USE AT YOUR OWN RISK.
func Slice(s string) (b []byte) {
	if len(s) == 0 {
		return []byte{}
	}
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pbytes.Data = pstring.Data
	pbytes.Len = pstring.Len
	pbytes.Cap = pstring.Len
	return
}

package xgen

import (
	"crypto/rand"
	"io"

	"github.com/VNTools1/xutils/xcrypto"
	"github.com/VNTools1/xutils/xencoding"
	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rs/xid"
	uuid "github.com/satori/go.uuid"
)

// GUID ...
func GUID() string {
	b := make([]byte, 48)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return ""
	}
	return xcrypto.Md5(xencoding.Base64Encode(string(b)))
}

// UUID ...
func UUID() string {
	return uuid.NewV4().String()
}

// XID ...
func XID() string {
	return xid.New().String()
}

// Nanoid ...
func Nanoid(l ...int) string {
	id, _ := nanoid.New(l...)
	return id
}

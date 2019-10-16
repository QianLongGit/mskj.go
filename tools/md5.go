package tools

import (
	"crypto/md5"
	"io"
	"fmt"
)

func Md5(raw string) string {
	h := md5.New()
	io.WriteString(h, raw)

	return fmt.Sprintf("%x", h.Sum(nil))
}

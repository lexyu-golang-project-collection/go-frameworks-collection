package main

import (
	"crypto/sha1"
	"encoding/base64"
	"io"
	"os"
)

func main() {
	h := sha1.New()
	io.WriteString(h, "YSBzYW1wbGUgMTYgYnl0ZQ==258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	encoder := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	encoder.Write(h.Sum(nil))
	// zkTGI6zVrOIDXiC4vnn1Rf37YFw=
	encoder.Close()
}

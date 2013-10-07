package main

import (
	"github.com/elegios/badbank/protocol"
	"net"
	"time"
)

const (
	timeout = time.Second
)

func handleBlobConn(blobConn net.Conn) (<-chan []string, chan<- string) {
	out := make(chan []string)
	in := make(chan string, 1)
	go func() {
		for {
			select {
			case lang := <-in:
				askForBlob(blobConn, lang)
			default:
				readBlob(blobConn, out)
			}
		}
	}()
	return out, in
}

func askForBlob(blobConn net.Conn, lang string) {
	m := &protocol.Message{protocol.Iam, 0, 0}
	m.SetASCII(lang)
	b, erp := m.Encode()
	d(erp)
	_, erp = blobConn.Write(b)
	d(erp)
}

func readBlob(blobConn net.Conn, out chan<- []string) {
	b := make([]byte, protocol.MaxBlobSize)

	blobConn.SetReadDeadline(time.Now().Add(timeout))
	_, err := blobConn.Read(b) //TODO, check correctness/error handling
	if err != nil {
		return
	}
	strings := protocol.DecodeBlob(b)
	if err != nil {
		return //TODO: print some error, the blob package was malformed
	}
	out <- strings
}

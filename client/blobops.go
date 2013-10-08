package main

import (
	"fmt"
	"github.com/elegios/badbank/protocol"
	"net"
	"time"
)

const (
	timeout = time.Second
)

var buffer []byte

func handleBlobConn(blobConn net.Conn) (<-chan []string, chan<- string) {
	buffer = make([]byte, protocol.MaxBlobSize)
	out := make(chan []string)
	in := make(chan string, 1)
	go func() {
		for {
			select {
			case lang := <-in:
				askForBlob(blobConn, lang)
			default:
				err := readBlob(blobConn, out)
				if err != nil {
					return
				}
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

func readBlob(blobConn net.Conn, out chan<- []string) (err error) {
	blobConn.SetReadDeadline(time.Now().Add(timeout))
	n, err := blobConn.Read(buffer) //TODO, check correctness/error handling
	if err != nil {
		return
	}
	strings := protocol.DecodeBlob(buffer[:n])
	if err != nil {
		fmt.Print("died")
		return //TODO: print some error, the blob package was malformed
	}
	out <- strings
	return
}

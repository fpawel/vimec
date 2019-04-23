package main

import (
	"encoding/binary"
	"io"
	"log"
	"time"
)

type Entity struct {
	ID            uint64
	Order         uint64
	Time          time.Time
	DocCode       byte
	ProductsCount uint64
	RouteSheet    string
}


func readBytes(x io.Reader, b []byte) {
	n, err := x.Read(b)
	if n != len(b) {
		log.Panicln(n, "!=len(b)==", len(b))
	}
	if err != nil && err != io.EOF {
		log.Panic(err)
	}
}

func readUInt64(x io.Reader) uint64 {
	b := make([]byte, 8)
	readBytes(x, b)
	return binary.BigEndian.Uint64(b)
}

func readString(x io.Reader) string {
	n := readUInt64(x)
	b := make([]byte, int(n))
	readBytes(x, b)
	return string(b)
}

func readEntity(k []byte, x io.Reader) (a Entity) {
	a.ID = binary.BigEndian.Uint64(k)
	a.Order = readUInt64(x)
	a.Time = time.Unix(0, int64(readUInt64(x)))
	b := []byte{0}
	readBytes(x, b)
	a.DocCode = b[0]
	a.ProductsCount = readUInt64(x)
	a.RouteSheet = readString(x)
	return a
}

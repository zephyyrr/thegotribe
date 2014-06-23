package thegotribe

import (
	"encoding/json"
	"io"
	"net"
)

type EyeTracker struct {
	enc json.Encoder
	dec json.Decoder
}

func New(conn io.ReadWriteCloser) *EyeTracker {
	return &EyeTracker{
		enc: json.NewEncoder(conn),
		dec: json.NewDecoder(conn),
	}
}

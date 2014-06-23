package thegotribe

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"io"
	"net"
)

type EyeTracker struct {
	enc      *json.Encoder
	dec      *json.Decoder
	incoming chan Response
	OnGaze   GazeFunc
}

func Create() (tracker *EyeTracker, err error) {
	conn, err := net.Dial("tcp", "localhost:6555")
	if err != nil {
		return
	}
	return New(conn), nil
}

func New(conn io.ReadWriteCloser) (tracker *EyeTracker) {
	tracker = &EyeTracker{
		enc: json.NewEncoder(conn),
		dec: json.NewDecoder(conn),
	}
	go tracker.readPackets()

	tracker.Set("version", ProtocolVersion)
	tracker.Set("push", true)

	// res := tracker.

	return
}

func (et *EyeTracker) readPackets() {
	var res Response
	var frame Frame
	mapDecoder, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{
			Result: &frame,
		})
	if err != nil {
		return
	}

	for {
		et.dec.Decode(res)
		if val, ok := res.Values["frame"]; ok {
			mapDecoder.Decode(val)
			et.OnGaze(frame)
		}
	}
}

func (et *EyeTracker) Set(attribute string, value interface{}) error {
	et.SetAll(map[string]interface{}{
		attribute: value,
	})
	return nil
}

func (et *EyeTracker) SetAll(attributes map[string]interface{}) error {
	et.enc.Encode(Request{
		Category: CategoryTracker,
		Request:  "set",
		Values:   attributes,
	})
	res := <-et.incoming
	if res.StatusCode != OK {
		return res.StatusCode
	}
	return nil
}

func (et *EyeTracker) Get(attribute string) (interface{}, error) {
	res, err := et.GetAll(attribute)
	return res[attribute], err
}

func (et *EyeTracker) GetAll(attributes ...string) (map[string]interface{}, error) {
	et.enc.Encode(Request{
		Category: CategoryTracker,
		Request:  "get",
		Values:   attributes,
	})

	return (<-et.incoming).Values, nil
}

type GazeFunc func(Frame)

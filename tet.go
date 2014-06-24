package thegotribe

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"io"
	"net"
	"time"
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
		enc:      json.NewEncoder(conn),
		dec:      json.NewDecoder(conn),
		incoming: make(chan Response, 10),
		OnGaze:   NullFunc,
	}
	go tracker.readPackets()

	tracker.Set("version", ProtocolVersion)
	tracker.Set("push", true)

	res, err := tracker.Get("heartbeatinterval")
	heartbeat := 250 * time.Millisecond
	if err == nil {
		heartbeat = res.(time.Duration) * time.Millisecond
	}
	go tracker.heartbeat(heartbeat)

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
		err := et.dec.Decode(res)
		if err != nil {
			close(et.incoming)
			return
		}
		if val, ok := res.Values["frame"]; ok {
			mapDecoder.Decode(val)
			et.OnGaze(frame)
		} else {
			et.incoming <- res
		}
	}
}

func (et *EyeTracker) heartbeat(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for _ = range ticker.C {
		et.enc.Encode(Request{
			Category: CategoryHeartbeat,
		})
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

	if err != nil {
		return nil, err
	} else {
		if val, ok := res[attribute]; !ok {
			return val, nil
		}
		return nil, errors.New(fmt.Sprintf("Unknown attribute \"%s\"", attribute))
	}
}

func (et *EyeTracker) GetAll(attributes ...string) (map[string]interface{}, error) {
	et.enc.Encode(Request{
		Category: CategoryTracker,
		Request:  "get",
		Values:   attributes,
	})

	result, ok := <-et.incoming
	if !ok {
		return nil, errors.New("End of Connection")
	}
	if result.StatusCode != OK {
		return result.Values, result.StatusCode
	}

	return result.Values, nil
}

type GazeFunc func(Frame)

func NullFunc(f Frame) {

}

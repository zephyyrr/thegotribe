package thegotribe

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"strings"
	"time"
)

var logger = log.New(os.Stdout, "[thegotribe] ", log.LstdFlags)

type EyeTracker struct {
	conn       io.ReadWriteCloser
	enc        *json.Encoder
	dec        *json.Decoder
	incoming   chan Response
	heartbeats chan Response
	OnGaze     GazeFunc
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
		conn:       conn,
		enc:        json.NewEncoder(conn),
		dec:        json.NewDecoder(conn),
		incoming:   make(chan Response, 10),
		heartbeats: make(chan Response),
		OnGaze:     NullFunc,
	}
	go tracker.readPackets()

	res, err := tracker.Get("heartbeatinterval")
	heartbeat := 250 * time.Millisecond
	if err == nil {
		requested := time.Duration(res.(float64)) * time.Millisecond
		if requested < heartbeat {
			heartbeat = requested / 2
		}
	}
	go tracker.heartbeat(heartbeat)

	tracker.Set("version", ProtocolVersion)
	tracker.Set("push", true)

	return
}

func (et *EyeTracker) Close() error {
	logger.Println("Closing tracker")
	return et.conn.Close()
}

func (et *EyeTracker) readPackets() {
	var res Response

	for {
		err := et.dec.Decode(&res)
		logger.Println("Raw:", res)
		if err != nil {
			logger.Println(err)
			close(et.incoming)
			close(et.heartbeats)
			return
		}

		if res.Category == CategoryHeartbeat {
			et.heartbeats <- res
		}

		if et.OnGaze == nil {
			logger.Println("OnGaze == nil")
		}

		if frame := res.Values.Frame; frame != nil && et.OnGaze != nil {
			logger.Println(frame)
			et.OnGaze(*frame)
			res.Values.Frame = nil // Handeled
		}
		if res.Values.StatusMessage != nil ||
			res.Values.Push != nil ||
			res.Values.HeartbeatInterval != nil ||
			res.Values.Version != nil ||
			res.Values.TrackerState != nil ||
			res.Values.Framerate != nil ||
			res.Values.ScreenIndex != nil ||
			res.Values.ScreenPsyW != nil ||
			res.Values.ScreenResW != nil ||
			res.Values.ScreenResH != nil ||
			res.Values.ScreenPsyH != nil {
			et.incoming <- res
		}
	}
}

func (et *EyeTracker) heartbeat(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for _ = range ticker.C {
		err := et.enc.Encode(Request{
			Category: CategoryHeartbeat,
		})
		if err != nil {
			ticker.Stop()
		}
		<-et.heartbeats
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
		if val := reflect.ValueOf(res).Elem().FieldByNameFunc(func(name string) bool {
			return strings.ToLower(name) == attribute
		}); !val.IsNil() {
			return val.Elem().Interface(), nil
		}
		return nil, errors.New(*res.StatusMessage)
	}
}

func (et *EyeTracker) GetAll(attributes ...string) (*Values, error) {
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
		return &result.Values, result.StatusCode
	}

	return &result.Values, nil
}

type GazeFunc func(Frame)

func NullFunc(f Frame) {

}

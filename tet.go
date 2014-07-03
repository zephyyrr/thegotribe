package thegotribe

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"reflect"
	"time"
)

const (
	DefaultAddr = "localhost:6555"
)

var logger = log.New(ioutil.Discard, "[thegotribe] ", log.LstdFlags)

type EyeTracker struct {
	conn io.ReadWriteCloser
	enc  *json.Encoder
	dec  *json.Decoder

	Frames     chan Frame
	heartbeats chan Response
	sets       chan Response
	interests  map[string][]chan<- interface{}
}

func Create() (tracker *EyeTracker, err error) {
	conn, err := net.Dial("tcp", DefaultAddr)
	if err != nil {
		return
	}
	return New(conn), nil
}

func New(conn io.ReadWriteCloser) (tracker *EyeTracker) {
	tracker = &EyeTracker{
		conn: conn,
		enc:  json.NewEncoder(conn),
		dec:  json.NewDecoder(conn),

		Frames:     make(chan Frame, 1),
		heartbeats: make(chan Response),
		sets:       make(chan Response),
		interests:  make(map[string][]chan<- interface{}),
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
			close(et.heartbeats)
			close(et.Frames)
			return
		}

		if res.StatusCode != OK { // Deals with bad status codes
			continue
		}

		if res.Category == CategoryHeartbeat {
			et.heartbeats <- res
			continue
		}

		if res.Category == CategoryCalibration {
			// Not yet supported.
			continue
		}

		// Only res.Category == CategoryTracker possible at this point
		if res.Category != CategoryTracker {
			panic("Unknown category " + res.Category)
		}

		if res.Request == RequestSet {
			et.sets <- res
			continue
		}

		// Only res.Request == RequestGet possible at this point
		if res.Request != RequestGet {
			panic("Unknown request method " + res.Request)
		}

		if frame := res.Values.Frame; frame != nil {
			select {
			case et.Frames <- *frame:
			default: //Throw it away. Client is busy.
			}
			res.Values.Frame = nil // Handeled
		}

		// Deal with every other value of tracker
		val, typ := reflect.ValueOf(res.Values), reflect.TypeOf(res.Values)
		for i := 0; i < val.NumField(); i++ {
			if field := val.Field(i); !field.IsNil() {
				name := typ.Field(i).Tag.Get("json") // Only if value is present
				if list, ok := et.interests[name]; ok && len(list) > 0 {
					for _, ch := range list { // Only if value is present and there are interested parties
						ch <- field.Elem().Interface() // Notify
						close(ch)                      // Close
					}
					et.interests[name] = list[0:0] //Clear interest list
				}
			}
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
	return et.SetAll(map[string]interface{}{
		attribute: value,
	})

}

func (et *EyeTracker) SetAll(attributes map[string]interface{}) error {
	et.enc.Encode(Request{
		Category: CategoryTracker,
		Request:  RequestSet,
		Values:   attributes,
	})
	res := <-et.sets
	if res.StatusCode != OK {
		return res.StatusCode
	}
	return nil
}

func (et *EyeTracker) Get(attribute string) (interface{}, error) {
	resch := make(chan interface{}, 1)
	et.interests[attribute] = append(et.interests[attribute], resch)

	err := et.enc.Encode(Request{
		Category: CategoryTracker,
		Request:  RequestGet,
		Values:   []string{attribute},
	})

	if err != nil {
		return nil, err
	}

	select {
	case res := <-resch:
		return res, nil
	case <-time.After(3 * time.Second):
		return nil, errors.New("Request timeout")
	}
}

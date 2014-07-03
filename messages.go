package thegotribe

import (
	"fmt"
	"reflect"
	"strings"
)

const ProtocolVersion = 1

type Request struct {
	Category Category    `json:"category"`
	Request  string      `json:"request"`
	Values   interface{} `json:"values"`
}

type Response struct {
	Category   Category `json:"category"`
	Request    string   `json:"request"`
	Values     Values   `json:"values"`
	StatusCode Status   `json:"statuscode"`
}

type Values struct {
	Frame         *Frame  `json:"frame"`
	StatusMessage *string `json:"statusmessage"`

	Push              *bool    `json:"push"`
	HeartbeatInterval *float64 `json:"heartbeatinterval"`
	Version           *float64 `json:"version"`
	TrackerState      *float64 `json:"trackerstate"`
	Framerate         *float64 `json:"framerate"`

	ScreenIndex *float64 `json:"screenindex"`
	ScreenResW  *float64 `json:"screenresw"`
	ScreenResH  *float64 `json:"screenresh"`
	ScreenPsyW  *float64 `json:"screenpsyw"`
	ScreenPsyH  *float64 `json:"screenpsyh"`

	IsCalibrated  *bool
	IsCalibrating *bool
	//Calibresult object
}

func (r Values) String() string {
	avail := make([]string, 0, 14)
	t := reflect.TypeOf(r)
	v := reflect.ValueOf(r)
	for i := 0; i < t.NumField(); i++ {
		if f := v.Field(i); !f.IsNil() {
			avail = append(avail, t.Field(i).Name)
		}
	}
	return `[` + strings.Join(avail, ", ") + `]`
}

type Category string

const (
	CategoryTracker     Category = "tracker"
	CategoryCalibration Category = "calibration"
	CategoryHeartbeat   Category = "heartbeat"
)

const (
	RequestGet = "get"
	RequestSet = "set"
)

type Status int

func (s Status) Error() string {
	switch s {
	case OK:
		return fmt.Sprintf("OK")
	case CalibrationChanged:
		return fmt.Sprintf("Calibration Changed")
	case DisplayChange:
		return fmt.Sprintf("Display Changed")
	case TrackerStateChange:
		return fmt.Sprintf("Tracker State Changed")
	}
	return fmt.Sprintf("Unknown error (%d)", int(s))
}

const (
	OK                 Status = 200
	CalibrationChanged        = 800
	DisplayChange             = 801
	TrackerStateChange        = 802
)

const (
	TrackerConnected = iota
	TrackerNotConnected
	TrackerConnectedBadFW
	TrackerConnectedNoUSB3
	TrackerConnectedNoStream
)

package thegotribe

import (
	"fmt"
)

const ProtocolVersion = 1

type Request struct {
	Category Category    `json:"category"`
	Request  string      `json:"request"`
	Values   interface{} `json:"values"`
}

type Response struct {
	Category   Category               `json:"category"`
	Request    string                 `json:"request"`
	Values     map[string]interface{} `json:"values"`
	StatusCode Status                 `json:"statuscode"`
}

type Category string

const (
	CategoryTracker     Category = "tracker"
	CategoryCalibration Category = "calibration"
	CategoryHeartbeat   Category = "heartbeat"
)

type Status int

func (s Status) Error() string {
	switch s {
	case OK:
		return fmt.Sprintf("%d: OK", s)
	case CalibrationChanged:
		return fmt.Sprintf("%d: Calibration Changed", s)
	case DisplayChange:
		return fmt.Sprintf("%d: Display Changed", s)
	case TrackerStateChange:
		return fmt.Sprintf("%d: Tracker State Changed", s)
	}
	return fmt.Sprintf("%d: Unknown error", s)
}

const (
	OK                 = 200
	CalibrationChanged = 800
	DisplayChange      = 801
	TrackerStateChange = 802
)

const (
	TrackerConnected = iota
	TrackerNotConnected
	TrackerConnectedBadFW
	TrackerConnectedNoUSB3
	TrackerConnectedNoStream
)

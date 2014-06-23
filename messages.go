package thegotribe

const ProtocolVersion = 1

type Request struct {
	Category Category    `json:"category"`
	Request  string      `json:"request"`
	Values   interface{} `json:"values"`
}

type Response struct {
	Request
	StatusCode int `json:"statuscode"`
}

type Category string

const (
	CategoryTracker     Category = "tracker"
	CategoryCalibration Category = "calibration"
	CategoryHeartbeat   Category = "heartbeat"
)

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

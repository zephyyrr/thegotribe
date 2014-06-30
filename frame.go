package thegotribe

import (
	"time"
)

type Frame struct {
	Timestamp TETTime `json:"timestamp",mapstructure:"timestamp"`
	Time      int64   `json:"time,mapstructure:"time""`
	Fix       bool    `json:"fix",mapstructure:"fix"`
	State     `json:"state",mapstructure:"state"`

	Raw     Point2D `json:"raw"`
	Average Point2D `json:"avg"`

	LeftEye  EyeData `json:"lefteye"`
	RightEye EyeData `json:"righteye"`
}

type Point2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type EyeData struct {
	Raw              Point2D `json:"raw"`
	Average          Point2D `json:"avg"`
	PupilSize        float64 `json:"psize"`
	PupilCoordinates Point2D `json:"pcenter"`
}

type State uint32

const (
	stateTrackingGaze = 0x1 << iota
	stateTrackingEyes
	stateTrackingPresence
	stateTrackingFail
	stateTrackingLost
)

func (s State) TrackingGaze() bool {
	return s&stateTrackingGaze != 0
}

func (s State) TrackingEyes() bool {
	return s&stateTrackingEyes != 0
}

func (s State) TrackingPresence() bool {
	return s&stateTrackingPresence != 0
}

func (s State) TrackingFail() bool {
	return s&stateTrackingFail != 0
}

func (s State) TrackingLost() bool {
	return s&stateTrackingLost != 0
}

type TETTime struct {
	time.Time
}

const (
	TETTimeLayout = "2006-01-02 15:04:05.999"
)

func (t TETTime) format() string {
	return t.Format(TETTimeLayout)
}

func (t TETTime) MarshalText() ([]byte, error) {
	return []byte(t.format()), nil
}

func (t TETTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.format() + `"`), nil
}

func (t *TETTime) UnmarshalJSON(data []byte) (err error) {
	t.Time, err = time.ParseInLocation(`"`+TETTimeLayout+`"`, string(data), time.Local)
	return
}

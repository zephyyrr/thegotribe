package thegotribe

import (
	"time"
)

type Frame struct {
	Timestamp time.Time `json:"timestamp"`
	Time      int64     `json:"time"`
	Fix       bool      `json:"fix"`
	State     `json:"state"`

	Raw     Point2D `json:"raw"`
	Average Point2D `json:"avg"`

	LeftEye  EyeData `json:"lefteye"`
	RightEye EyeData `json:"righteye"`
}

type Point2D struct {
	X int `json:"x"`
	Y int `json:"y"`
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

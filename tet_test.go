package thegotribe

import (
	"testing"
)

func TestConnect(t *testing.T) {
	et, err := Create()
	if err != nil {
		t.Fatal("Failed to create tracker:", err)
	}
	if res, err := et.Get("push"); err != nil || res.(bool) != true {
		if err != nil {
			t.Error(err)
		} else {
			t.Error("Tracker not in push-mode.")
		}
	}
}

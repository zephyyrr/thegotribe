package thegotribe

import (
	"testing"
)

func TestConnect(t *testing.T) {
	_, err := Create()
	if err != nil {
		t.Fatal("Failed to create tracker:", err)
	}
	//t.Log(et.Get("push"))
}

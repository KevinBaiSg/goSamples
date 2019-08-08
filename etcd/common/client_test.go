package common

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	if _, err := NewClient(); err != nil {
		t.Errorf("new client error %e", err)
		return
	}
	t.Logf("new client success")
}

func ExampleNewClient() {
	client, e := NewClient()
	if e != nil {
		return
	}
	defer client.Close()
}

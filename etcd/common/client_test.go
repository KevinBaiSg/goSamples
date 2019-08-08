package common

import (
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestNewClient(t *testing.T) {
	dir, err := filepath.Abs("./")
	if err != nil {
		t.Errorf("filepath directory error %e", err)
		return
	}
	viper.AddConfigPath(dir)

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

	// Output: dsfa
}

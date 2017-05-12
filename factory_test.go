package urlshortner

import (
	"testing"
)

var factory *Factory
var key = "12345"
var long = "google.com"

func TestMain(t *testing.M) {
	factory = NewFactory("127.0.0.1", 5984, "tinyurl", "test_design")
}

func TestAdd(t *testing.T) {
	err := factory.AddURL(key, long)
	if err != nil {
		t.Error(err)
	}
}

func TestGetURL(t *testing.T) {
	_, err := factory.GetURL(key)
	if err != nil {
		t.Error(err)
	}
}

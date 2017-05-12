package random

import "testing"

func TestRandom(t *testing.T) {
	rand, err := NewRandom()
	if err != nil {
		t.Error(err)
	}

	_, err = rand.GetRandomUrl(9)
	if err != nil {
		t.Error(err)
	}
}

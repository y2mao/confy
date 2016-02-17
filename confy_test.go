package confy

import (
	"fmt"
	"testing"
)

func TestConfy(t *testing.T) {
	Define("http.host", "127.0.0.1")
	Define("http.port", 2009)
	Define("http.auth.enabled", true)

	Ready()

	assertEqual(t, Text("http.host"), "127.0.0.1")
	assertEqual(t, Int("http.port"), 2009)
	assertEqual(t, Bool("http.auth.enabled"), true)
}

func assertEqual(t *testing.T, v1 interface{}, v2 interface{}) {
	s1 := fmt.Sprintf("%v", v1)
	s2 := fmt.Sprintf("%v", v2)
	if s1 != s2 {
		t.Errorf("Not Equal: %#v != %#v", v1, v2)
	}
}

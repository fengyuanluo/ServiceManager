package manager

import "testing"

func TestParsePM2List(t *testing.T) {
	out := "app online"
	items, _ := ParsePM2List(out)
	if len(items) != 1 || items[0].Name != "app" {
		t.Fatalf("parse failed")
	}
}

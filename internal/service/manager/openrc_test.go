package manager

import "testing"

func TestParseOpenRCList(t *testing.T) {
	out := "ssh [ started ]"
	items, _ := ParseOpenRCList(out)
	if len(items) != 1 || items[0].Name != "ssh" {
		t.Fatalf("parse failed")
	}
}

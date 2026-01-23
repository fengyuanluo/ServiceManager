package manager

import "testing"

func TestParseOpenRCList(t *testing.T) {
	out := "ssh [ started ]\ncron [ stopped ]"
	items, _ := ParseOpenRCList(out)
	if len(items) != 2 || items[0].Name != "ssh" || items[1].Name != "cron" {
		t.Fatalf("parse failed")
	}
}

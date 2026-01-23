package manager

import "testing"

func TestParsePM2List(t *testing.T) {
	out := `[{"name":"app","pid":123,"pm2_env":{"status":"online"}},{"name":"job","pid":0,"pm2_env":{"status":"stopped"}}]`
	items, _ := ParsePM2List(out)
	if len(items) != 2 || items[0].Name != "app" || items[1].Name != "job" {
		t.Fatalf("parse failed")
	}
}

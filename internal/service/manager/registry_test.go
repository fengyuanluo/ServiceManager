package manager

import "testing"

func TestDetectManagers(t *testing.T) {
	env := EnvInfo{Distro: "ubuntu", Init: "systemd"}
	list := DetectManagers(env)
	if len(list) == 0 {
		t.Fatalf("expected managers")
	}
}

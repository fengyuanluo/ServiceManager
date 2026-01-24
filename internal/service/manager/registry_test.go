package manager

import (
	"reflect"
	"testing"
)

func TestDetectManagers(t *testing.T) {
	cases := []struct {
		env  EnvInfo
		want []string
	}{
		{env: EnvInfo{Distro: "ubuntu", Init: "systemd"}, want: []string{"systemd"}},
		{env: EnvInfo{Distro: "alpine", Init: "openrc"}, want: []string{"openrc"}},
		{env: EnvInfo{Distro: "debian", Init: "init"}, want: []string{"sysv"}},
		{env: EnvInfo{Distro: "unknown", Init: "unknown"}, want: []string{}},
	}
	for _, item := range cases {
		got := DetectManagers(item.env)
		if !reflect.DeepEqual(got, item.want) {
			t.Fatalf("unexpected managers for %+v: %v", item.env, got)
		}
	}
}

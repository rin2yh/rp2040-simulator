//go:build !tinygo

package bridge

import (
	"slices"
	"testing"
)

func TestEmulatorEnvironmentReplacesInternalSettings(t *testing.T) {
	env := emulatorEnvironment([]string{
		"APP_SETTING=value",
		EmulatorProcess + "=stale",
		EmulatorBoard + "=stale",
	}, "conf2025badge")
	want := []string{
		"APP_SETTING=value",
		EmulatorProcess + "=1",
		EmulatorBoard + "=conf2025badge",
	}
	if !slices.Equal(env, want) {
		t.Fatalf("environment = %q, want %q", env, want)
	}
}

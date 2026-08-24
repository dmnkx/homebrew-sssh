package cli

import (
	"strings"
	"testing"
)

func TestShow(t *testing.T) {
	p := writeConfig(t, `Host prod
    HostName 10.0.0.1
    User ubuntu
`)
	out, err := runCLI(t, "--config", p, "show", "prod")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Host prod") || !strings.Contains(out, "HostName 10.0.0.1") {
		t.Fatal(out)
	}
}

func TestShowUnknown(t *testing.T) {
	p := writeConfig(t, "Host a\n    HostName 1.1.1.1\n")
	if _, err := runCLI(t, "--config", p, "show", "missing"); err == nil {
		t.Fatal("expected unknown host")
	}
}

func TestShowRequiresAlias(t *testing.T) {
	p := writeConfig(t, "Host a\n    HostName 1.1.1.1\n")
	if _, err := runCLI(t, "--config", p, "show"); err == nil {
		t.Fatal("show without alias should fail")
	}
}

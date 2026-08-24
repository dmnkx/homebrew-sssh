package cli

import (
	"strings"
	"testing"
)

func TestConnectUnknown(t *testing.T) {
	p := writeConfig(t, "Host a\n    HostName 1.1.1.1\n")
	if _, err := runCLI(t, "--config", p, "connect", "nope"); err == nil {
		t.Fatal("expected unknown host")
	}
}

func TestConnectPrintCmd(t *testing.T) {
	putFakeSSH(t)
	p := writeConfig(t, "Host prod\n    HostName 10.0.0.1\n")
	out, err := runCLI(t, "--print-cmd", "--config", p, "connect", "prod")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "prod") {
		t.Fatal(out)
	}
}

func TestConnectRequiresAlias(t *testing.T) {
	p := writeConfig(t, "Host a\n    HostName 1.1.1.1\n")
	if _, err := runCLI(t, "--config", p, "connect"); err == nil {
		t.Fatal("connect without alias should fail")
	}
}

package cli

import (
	"strings"
	"testing"
)

func TestList(t *testing.T) {
	p := writeConfig(t, `Host prod
    HostName 10.0.0.1
    User ubuntu
    Port 2222
    IdentityFile ~/.ssh/id_ed25519
`)
	out, err := runCLI(t, "--config", p, "list")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "prod") || !strings.Contains(out, "ubuntu@10.0.0.1:2222") {
		t.Fatal(out)
	}
	if !strings.Contains(out, "id_ed25519") {
		t.Fatal(out)
	}
}

func TestListDefaultPortAndEmptyFields(t *testing.T) {
	p := writeConfig(t, "Host bare\n")
	out, err := runCLI(t, "--config", p, "list")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "bare") || !strings.Contains(out, "-@-:22") {
		t.Fatal(out)
	}
}

func TestListEmptyConfig(t *testing.T) {
	p := writeConfig(t, "")
	out, err := runCLI(t, "--config", p, "list")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(out) != "" {
		t.Fatal(out)
	}
}

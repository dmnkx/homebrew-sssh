package cli

import (
	"testing"

	"sssh/internal/sshcfg"
)

func TestAdd(t *testing.T) {
	p := writeConfig(t, "")
	_, err := runCLI(t, "--config", p, "add", "dev",
		"--host", "1.2.3.4",
		"--user", "ec2",
		"--port", "22",
		"--key", "~/.ssh/id_ed25519",
		"--jump", "bastion",
	)
	if err != nil {
		t.Fatal(err)
	}
	f, err := sshcfg.Load(p)
	if err != nil {
		t.Fatal(err)
	}
	h, ok := f.Find("dev")
	if !ok {
		t.Fatal("missing host")
	}
	if h.HostName != "1.2.3.4" || h.User != "ec2" || h.Port != "22" {
		t.Fatalf("%#v", h)
	}
	if h.IdentityFile != "~/.ssh/id_ed25519" || h.ProxyJump != "bastion" {
		t.Fatalf("%#v", h)
	}
}

func TestAddDuplicate(t *testing.T) {
	p := writeConfig(t, "")
	_, err := runCLI(t, "--config", p, "add", "dev", "--host", "1.2.3.4")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runCLI(t, "--config", p, "add", "dev", "--host", "9.9.9.9"); err == nil {
		t.Fatal("duplicate add should fail")
	}
}

func TestAddRequiresHostFlag(t *testing.T) {
	p := writeConfig(t, "")
	if _, err := runCLI(t, "--config", p, "add", "x"); err == nil {
		t.Fatal("add without --host should fail")
	}
}

func TestAddRequiresAlias(t *testing.T) {
	p := writeConfig(t, "")
	if _, err := runCLI(t, "--config", p, "add", "--host", "1.2.3.4"); err == nil {
		t.Fatal("add without alias should fail")
	}
}

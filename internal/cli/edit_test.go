package cli

import (
	"testing"

	"sssh/internal/sshcfg"
)

func TestEdit(t *testing.T) {
	p := writeConfig(t, "")
	_, err := runCLI(t, "--config", p, "add", "dev", "--host", "1.2.3.4", "--user", "ec2")
	if err != nil {
		t.Fatal(err)
	}
	_, err = runCLI(t, "--config", p, "edit", "dev", "--user", "ubuntu", "--port", "2222", "--jump", "bastion", "--key", "/tmp/key")
	if err != nil {
		t.Fatal(err)
	}
	f, err := sshcfg.Load(p)
	if err != nil {
		t.Fatal(err)
	}
	h, _ := f.Find("dev")
	if h.User != "ubuntu" || h.HostName != "1.2.3.4" {
		t.Fatalf("partial edit: %#v", h)
	}
	if h.Port != "2222" || h.ProxyJump != "bastion" || h.IdentityFile != "/tmp/key" {
		t.Fatalf("%#v", h)
	}
}

func TestEditHostName(t *testing.T) {
	p := writeConfig(t, "")
	_, err := runCLI(t, "--config", p, "add", "dev", "--host", "old.example")
	if err != nil {
		t.Fatal(err)
	}
	_, err = runCLI(t, "--config", p, "edit", "dev", "--host", "new.example")
	if err != nil {
		t.Fatal(err)
	}
	f, err := sshcfg.Load(p)
	if err != nil {
		t.Fatal(err)
	}
	h, _ := f.Find("dev")
	if h.HostName != "new.example" {
		t.Fatalf("%#v", h)
	}
}

func TestEditUnknown(t *testing.T) {
	p := writeConfig(t, "")
	if _, err := runCLI(t, "--config", p, "edit", "nope", "--user", "x"); err == nil {
		t.Fatal("expected unknown host")
	}
}

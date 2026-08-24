package cli

import (
	"testing"

	"sssh/internal/sshcfg"
)

func TestRm(t *testing.T) {
	p := writeConfig(t, "")
	_, err := runCLI(t, "--config", p, "add", "dev", "--host", "1.2.3.4")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runCLI(t, "--config", p, "rm", "dev"); err != nil {
		t.Fatal(err)
	}
	f, err := sshcfg.Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := f.Find("dev"); ok {
		t.Fatal("host still exists")
	}
}

func TestRmMissing(t *testing.T) {
	p := writeConfig(t, "")
	if _, err := runCLI(t, "--config", p, "rm", "dev"); err == nil {
		t.Fatal("rm missing should fail")
	}
}

func TestRmRequiresAlias(t *testing.T) {
	p := writeConfig(t, "")
	if _, err := runCLI(t, "--config", p, "rm"); err == nil {
		t.Fatal("rm without alias should fail")
	}
}

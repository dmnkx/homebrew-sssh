package connect

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestErrString(t *testing.T) {
	if got := errString(nil); got != "" {
		t.Fatalf("nil -> %q", got)
	}
	if got := errString(os.ErrNotExist); got == "" {
		t.Fatal("expected message")
	}
}

func TestSSHPathMissing(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	if _, err := SSHPath(); err == nil {
		t.Fatal("expected error when ssh is not on PATH")
	}
	if _, err := PrintCmd("prod"); err == nil {
		t.Fatal("PrintCmd should fail without ssh")
	}
}

func TestExecMissingSSH(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	if err := Exec("prod"); err == nil {
		t.Fatal("Exec should fail when ssh is not on PATH")
	}
}

func TestSSHPathFound(t *testing.T) {
	dir := t.TempDir()
	fake := filepath.Join(dir, "ssh")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	p, err := SSHPath()
	if err != nil {
		t.Fatal(err)
	}
	if p != fake {
		t.Fatalf("got %q want %q", p, fake)
	}
}

func TestPrintCmd(t *testing.T) {
	dir := t.TempDir()
	fake := filepath.Join(dir, "ssh")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	got, err := PrintCmd("prod")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(got, " ssh prod") && got != fake+" prod" {
		t.Fatalf("got %q", got)
	}
	if !strings.Contains(got, "prod") {
		t.Fatalf("got %q", got)
	}
}

package cli

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sssh/internal/sshcfg"
)

func TestEmpty(t *testing.T) {
	if empty("  ", "-") != "-" {
		t.Fatal("blank should use default")
	}
	if empty("ubuntu", "-") != "ubuntu" {
		t.Fatal("non-empty should keep value")
	}
}

func TestNewRootRegistersCommands(t *testing.T) {
	root := NewRoot()
	want := []string{"list", "add", "edit", "rm", "show", "connect"}
	for _, name := range want {
		if root.Commands() == nil {
			t.Fatal("no commands")
		}
		found := false
		for _, c := range root.Commands() {
			if c.Name() == name {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("missing command %s", name)
		}
	}
}

func TestTooManyArgs(t *testing.T) {
	p := writeConfig(t, "Host a\n    HostName 1.1.1.1\n")
	if _, err := runCLI(t, "--config", p, "a", "b"); err == nil {
		t.Fatal("two positional args should fail")
	}
}

func TestLoadExplicitPath(t *testing.T) {
	p := writeConfig(t, "Host a\n    HostName 1.1.1.1\n")
	f, err := load(p)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := f.Find("a"); !ok {
		t.Fatal("expected host a")
	}
}

func TestLoadMissingFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "missing")
	f, err := load(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(f.SelectableHosts()) != 0 {
		t.Fatalf("%#v", f.SelectableHosts())
	}
}

func TestConnectAliasUnknown(t *testing.T) {
	f := &sshcfg.File{}
	if err := connectAlias(f, "nope", true); err == nil {
		t.Fatal("expected unknown host")
	}
}

func TestConnectAliasPrint(t *testing.T) {
	putFakeSSH(t)
	p := writeConfig(t, "Host prod\n    HostName 10.0.0.1\n")
	f, err := sshcfg.Load(p)
	if err != nil {
		t.Fatal(err)
	}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	os.Stdout = w
	err = connectAlias(f, "prod", true)
	_ = w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatal(err)
	}
	got, readErr := io.ReadAll(r)
	_ = r.Close()
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(got), "prod") {
		t.Fatal(string(got))
	}
}

func TestLoadDefaultPath(t *testing.T) {
	if _, err := load(""); err != nil {
		t.Fatal(err)
	}
}

func TestRootPrintCmd(t *testing.T) {
	putFakeSSH(t)
	p := writeConfig(t, "Host prod\n    HostName 10.0.0.1\n")
	out, err := runCLI(t, "--config", p, "--print-cmd", "prod")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "prod") {
		t.Fatal(out)
	}
}

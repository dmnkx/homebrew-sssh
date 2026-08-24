package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func writeConfig(t *testing.T, body string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "config")
	if err := os.WriteFile(p, []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func runCLI(t *testing.T, args ...string) (string, error) {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = old })

	var errBuf bytes.Buffer
	root := NewRoot()
	root.SetOut(&errBuf)
	root.SetErr(&errBuf)
	root.SilenceUsage = true
	root.SetArgs(args)
	runErr := root.Execute()

	_ = w.Close()
	os.Stdout = old
	var out bytes.Buffer
	_, _ = out.ReadFrom(r)
	_ = r.Close()
	return out.String() + errBuf.String(), runErr
}

func putFakeSSH(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "ssh")
	if err := os.WriteFile(p, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
}

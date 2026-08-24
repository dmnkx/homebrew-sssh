package sshcfg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseAndUpsert(t *testing.T) {
	raw := `# comment
Host *
    IdentitiesOnly yes

Host prod staging
    HostName 10.0.0.1
    User ubuntu
    IdentityFile ~/.ssh/id_ed25519
`
	f := &File{Path: "mem", Blocks: parseBlocks(raw)}
	hosts := f.SelectableHosts()
	if len(hosts) != 2 {
		t.Fatalf("want 2 selectable, got %d", len(hosts))
	}
	if err := f.Upsert(Host{Alias: "dev", HostName: "1.2.3.4", User: "ec2"}); err != nil {
		t.Fatal(err)
	}
	out := f.Serialize()
	if !strings.Contains(out, "Host dev") {
		t.Fatalf("missing Host dev: %s", out)
	}
	if !strings.Contains(out, "Host *") {
		t.Fatalf("lost wildcard: %s", out)
	}
}

func TestSaveRoundTrip(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "config")
	f := &File{Path: p}
	if err := f.Upsert(Host{Alias: "a", HostName: "h", User: "u", Port: "22"}); err != nil {
		t.Fatal(err)
	}
	if err := f.Save(); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "Host a") {
		t.Fatal(string(got))
	}
	f2, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	h, ok := f2.Find("a")
	if !ok || h.User != "u" {
		t.Fatalf("%v %v", ok, h)
	}
}

func TestLoadMissingFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "missing")
	f, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if f.Path != p {
		t.Fatalf("path = %s", f.Path)
	}
	if len(f.SelectableHosts()) != 0 {
		t.Fatalf("want empty, got %#v", f.SelectableHosts())
	}
}

func TestFindUnknown(t *testing.T) {
	f := &File{Blocks: parseBlocks("Host a\n    HostName 1.1.1.1\n")}
	if _, ok := f.Find("nope"); ok {
		t.Fatal("expected miss")
	}
}

func TestUpsertRejectsEmptyAndWildcard(t *testing.T) {
	f := &File{}
	if err := f.Upsert(Host{}); err == nil {
		t.Fatal("empty alias should fail")
	}
	if err := f.Upsert(Host{Alias: "*"}); err == nil {
		t.Fatal("wildcard alias should fail")
	}
}

func TestUpsertReplacesExisting(t *testing.T) {
	f := &File{}
	if err := f.Upsert(Host{Alias: "a", HostName: "old"}); err != nil {
		t.Fatal(err)
	}
	if err := f.Upsert(Host{Alias: "a", HostName: "new", User: "me"}); err != nil {
		t.Fatal(err)
	}
	h, ok := f.Find("a")
	if !ok || h.HostName != "new" || h.User != "me" {
		t.Fatalf("%v %#v", ok, h)
	}
}

func TestRemoveSingleAndMissing(t *testing.T) {
	f := &File{}
	_ = f.Upsert(Host{Alias: "a", HostName: "h"})
	if err := f.Remove("a"); err != nil {
		t.Fatal(err)
	}
	if _, ok := f.Find("a"); ok {
		t.Fatal("still present")
	}
	if err := f.Remove("a"); err == nil {
		t.Fatal("second remove should fail")
	}
}

func TestRemoveOneAliasFromSharedBlock(t *testing.T) {
	f := &File{Blocks: parseBlocks("Host prod staging\n    HostName 10.0.0.1\n")}
	if err := f.Remove("staging"); err != nil {
		t.Fatal(err)
	}
	if _, ok := f.Find("staging"); ok {
		t.Fatal("staging should be gone")
	}
	h, ok := f.Find("prod")
	if !ok || h.HostName != "10.0.0.1" {
		t.Fatalf("%v %#v", ok, h)
	}
	if !strings.Contains(f.Serialize(), "Host prod") {
		t.Fatal(f.Serialize())
	}
	if strings.Contains(f.Serialize(), "staging") {
		t.Fatal(f.Serialize())
	}
}

func TestParseProxyJumpPortIdentityAndExtra(t *testing.T) {
	raw := `Host jumpbox
    HostName 8.8.8.8
    User root
    Port 2222
    IdentityFile ~/.ssh/id_rsa
    ProxyJump bastion
    ForwardAgent yes
`
	f := &File{Blocks: parseBlocks(raw)}
	h, ok := f.Find("jumpbox")
	if !ok {
		t.Fatal("missing host")
	}
	if h.Port != "2222" || h.ProxyJump != "bastion" || h.IdentityFile != "~/.ssh/id_rsa" {
		t.Fatalf("%#v", h)
	}
	if len(h.Extra) != 1 || h.Extra[0].Key != "ForwardAgent" || h.Extra[0].Value != "yes" {
		t.Fatalf("extra: %#v", h.Extra)
	}
	block := h.BlockString()
	if !strings.Contains(block, "ProxyJump bastion") || !strings.Contains(block, "ForwardAgent yes") {
		t.Fatal(block)
	}
}

func TestWildcardHostNotSelectable(t *testing.T) {
	f := &File{Blocks: parseBlocks("Host *\n    IdentitiesOnly yes\n")}
	if n := len(f.SelectableHosts()); n != 0 {
		t.Fatalf("got %d", n)
	}
}

func TestMatchBlockIgnored(t *testing.T) {
	raw := `Match host *.internal
    User git

Host real
    HostName 1.2.3.4
`
	f := &File{Blocks: parseBlocks(raw)}
	hosts := f.SelectableHosts()
	if len(hosts) != 1 || hosts[0].Alias != "real" {
		t.Fatalf("%#v", hosts)
	}
}

func TestSaveWritesBackup(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "config")
	if err := os.WriteFile(p, []byte("Host old\n    HostName 0.0.0.0\n"), 0600); err != nil {
		t.Fatal(err)
	}
	f, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Upsert(Host{Alias: "new", HostName: "1.1.1.1"}); err != nil {
		t.Fatal(err)
	}
	if err := f.Save(); err != nil {
		t.Fatal(err)
	}
	bak, err := os.ReadFile(p + ".bak")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(bak), "Host old") {
		t.Fatal(string(bak))
	}
}

func TestDefaultPath(t *testing.T) {
	p, err := DefaultPath()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(p, filepath.Join(".ssh", "config")) {
		t.Fatalf("unexpected path: %s", p)
	}
}

func TestSplitKey(t *testing.T) {
	k, rest := splitKey("HostName 10.0.0.1")
	if k != "HostName" || rest != "10.0.0.1" {
		t.Fatalf("%q %q", k, rest)
	}
	k, rest = splitKey("# comment")
	if k != "" || rest != "" {
		t.Fatalf("comment: %q %q", k, rest)
	}
	k, rest = splitKey("")
	if k != "" || rest != "" {
		t.Fatalf("empty: %q %q", k, rest)
	}
}

func TestIsWildcardAndFirstConcreteAlias(t *testing.T) {
	if !isWildcard("*") || !isWildcard("*.example") || !isWildcard("host?") {
		t.Fatal("expected wildcards")
	}
	if isWildcard("prod") {
		t.Fatal("prod is not a wildcard")
	}
	if got := firstConcreteAlias([]string{"*", "prod"}); got != "prod" {
		t.Fatalf("got %q", got)
	}
	if got := firstConcreteAlias([]string{"*", "?"}); got != "" {
		t.Fatalf("got %q", got)
	}
}

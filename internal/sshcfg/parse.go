package sshcfg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type Host struct {
	Alias        string
	HostName     string
	User         string
	Port         string
	IdentityFile string
	ProxyJump    string
	Extra        []KV
}

type KV struct {
	Key   string
	Value string
}

type File struct {
	Path   string
	Blocks []Block
}

type Block struct {
	Kind    string // "host" or "other"
	Aliases []string
	Lines   []string
	Host    Host
}

func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".ssh", "config"), nil
}

func Load(path string) (*File, error) {
	f := &File{Path: path}
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return f, nil
		}
		return nil, err
	}
	f.Blocks = parseBlocks(string(raw))
	return f, nil
}

func parseBlocks(s string) []Block {
	sc := bufio.NewScanner(strings.NewReader(s))
	var blocks []Block
	var cur *Block
	flush := func() {
		if cur == nil {
			return
		}
		if cur.Kind == "host" {
			cur.Host = hostFromBlock(*cur)
		}
		blocks = append(blocks, *cur)
		cur = nil
	}
	for sc.Scan() {
		line := sc.Text()
		trim := strings.TrimSpace(line)
		key, rest := splitKey(trim)
		if strings.EqualFold(key, "Host") {
			flush()
			aliases := strings.Fields(rest)
			cur = &Block{Kind: "host", Aliases: aliases, Lines: []string{line}}
			continue
		}
		if strings.EqualFold(key, "Match") {
			flush()
			cur = &Block{Kind: "other", Lines: []string{line}}
			continue
		}
		if cur == nil {
			cur = &Block{Kind: "other"}
		}
		cur.Lines = append(cur.Lines, line)
	}
	flush()
	return blocks
}

func splitKey(trim string) (key, rest string) {
	if trim == "" || strings.HasPrefix(trim, "#") {
		return "", ""
	}
	i := 0
	for i < len(trim) && !unicode.IsSpace(rune(trim[i])) {
		i++
	}
	if i == 0 {
		return "", ""
	}
	return trim[:i], strings.TrimSpace(trim[i:])
}

func hostFromBlock(b Block) Host {
	h := Host{Alias: firstConcreteAlias(b.Aliases)}
	if h.Alias == "" && len(b.Aliases) > 0 {
		h.Alias = b.Aliases[0]
	}
	for _, line := range b.Lines[1:] {
		trim := strings.TrimSpace(line)
		key, rest := splitKey(trim)
		if key == "" {
			continue
		}
		switch strings.ToLower(key) {
		case "hostname":
			h.HostName = rest
		case "user":
			h.User = rest
		case "port":
			h.Port = rest
		case "identityfile":
			h.IdentityFile = rest
		case "proxyjump":
			h.ProxyJump = rest
		default:
			h.Extra = append(h.Extra, KV{Key: key, Value: rest})
		}
	}
	return h
}

func firstConcreteAlias(aliases []string) string {
	for _, a := range aliases {
		if !isWildcard(a) {
			return a
		}
	}
	return ""
}

func isWildcard(a string) bool {
	return strings.ContainsAny(a, "*?!")
}

func (f *File) SelectableHosts() []Host {
	var out []Host
	seen := map[string]bool{}
	for _, b := range f.Blocks {
		if b.Kind != "host" {
			continue
		}
		for _, a := range b.Aliases {
			if isWildcard(a) || seen[a] {
				continue
			}
			seen[a] = true
			h := b.Host
			h.Alias = a
			out = append(out, h)
		}
	}
	return out
}

func (f *File) Find(alias string) (Host, bool) {
	for _, h := range f.SelectableHosts() {
		if h.Alias == alias {
			return h, true
		}
	}
	return Host{}, false
}

func (f *File) Upsert(h Host) error {
	if h.Alias == "" {
		return fmt.Errorf("empty alias")
	}
	if isWildcard(h.Alias) {
		return fmt.Errorf("wildcard alias not allowed: %s", h.Alias)
	}
	idx := f.hostIndex(h.Alias)
	block := blockFromHost(h)
	if idx >= 0 {
		f.Blocks[idx] = block
		return nil
	}
	f.Blocks = append(f.Blocks, block)
	return nil
}

func (f *File) Remove(alias string) error {
	idx := f.hostIndex(alias)
	if idx < 0 {
		return fmt.Errorf("host not found: %s", alias)
	}
	b := f.Blocks[idx]
	if len(b.Aliases) == 1 {
		f.Blocks = append(f.Blocks[:idx], f.Blocks[idx+1:]...)
		return nil
	}
	var next []string
	for _, a := range b.Aliases {
		if a != alias {
			next = append(next, a)
		}
	}
	b.Aliases = next
	if len(b.Lines) > 0 {
		b.Lines[0] = "Host " + strings.Join(next, " ")
	}
	b.Host.Alias = firstConcreteAlias(next)
	f.Blocks[idx] = b
	return nil
}

func (f *File) hostIndex(alias string) int {
	for i, b := range f.Blocks {
		if b.Kind != "host" {
			continue
		}
		for _, a := range b.Aliases {
			if a == alias {
				return i
			}
		}
	}
	return -1
}

func blockFromHost(h Host) Block {
	lines := []string{"Host " + h.Alias}
	add := func(k, v string) {
		if v != "" {
			lines = append(lines, "    "+k+" "+v)
		}
	}
	add("HostName", h.HostName)
	add("User", h.User)
	add("Port", h.Port)
	add("IdentityFile", h.IdentityFile)
	add("ProxyJump", h.ProxyJump)
	for _, kv := range h.Extra {
		add(kv.Key, kv.Value)
	}
	b := Block{Kind: "host", Aliases: []string{h.Alias}, Lines: lines}
	b.Host = h
	return b
}

func (f *File) Serialize() string {
	var b strings.Builder
	for i, block := range f.Blocks {
		for _, line := range block.Lines {
			b.WriteString(line)
			b.WriteByte('\n')
		}
		if i < len(f.Blocks)-1 {
			last := ""
			if len(block.Lines) > 0 {
				last = block.Lines[len(block.Lines)-1]
			}
			if strings.TrimSpace(last) != "" {
				nextFirst := ""
				if len(f.Blocks[i+1].Lines) > 0 {
					nextFirst = strings.TrimSpace(f.Blocks[i+1].Lines[0])
				}
				if nextFirst != "" {
					b.WriteByte('\n')
				}
			}
		}
	}
	return b.String()
}

func (f *File) Save() error {
	dir := filepath.Dir(f.Path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	if _, err := os.Stat(f.Path); err == nil {
		bak := f.Path + ".bak"
		data, err := os.ReadFile(f.Path)
		if err != nil {
			return err
		}
		if err := os.WriteFile(bak, data, 0600); err != nil {
			return err
		}
	}
	return os.WriteFile(f.Path, []byte(f.Serialize()), 0600)
}

func (h Host) BlockString() string {
	return strings.Join(blockFromHost(h).Lines, "\n")
}

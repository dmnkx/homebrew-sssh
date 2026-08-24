package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"sssh/internal/sshcfg"
)

func sampleHosts() []sshcfg.Host {
	return []sshcfg.Host{
		{Alias: "prod", HostName: "10.0.0.1", User: "ubuntu", Port: "22"},
		{Alias: "dev", HostName: "10.0.0.2", User: "ec2", Port: "2222"},
	}
}

func TestRunEmptyHosts(t *testing.T) {
	_, err := Run(nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFiltered(t *testing.T) {
	m := model{hosts: sampleHosts()}
	if n := len(m.filtered()); n != 2 {
		t.Fatalf("empty filter: %d", n)
	}
	m.filter = "PROD"
	got := m.filtered()
	if len(got) != 1 || got[0].Alias != "prod" {
		t.Fatalf("%#v", got)
	}
	m.filter = "ec2"
	got = m.filtered()
	if len(got) != 1 || got[0].Alias != "dev" {
		t.Fatalf("%#v", got)
	}
	m.filter = "zzz"
	if n := len(m.filtered()); n != 0 {
		t.Fatalf("want 0, got %d", n)
	}
}

func TestUpdateEscQuits(t *testing.T) {
	m := model{hosts: sampleHosts()}
	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	nm := next.(model)
	if !nm.quit || cmd == nil {
		t.Fatalf("quit=%v cmd=%v", nm.quit, cmd)
	}
}

func TestUpdateEnterSelects(t *testing.T) {
	m := model{hosts: sampleHosts(), cursor: 1}
	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	nm := next.(model)
	if nm.choice == nil || nm.choice.Alias != "dev" || cmd == nil {
		t.Fatalf("%#v", nm.choice)
	}
}

func TestUpdateEnterEmptyFilterList(t *testing.T) {
	m := model{hosts: sampleHosts(), filter: "nomatch"}
	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	nm := next.(model)
	if nm.choice != nil || cmd != nil {
		t.Fatal("should stay idle")
	}
}

func TestUpdateFilterAndView(t *testing.T) {
	m := model{hosts: sampleHosts()}
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("ec2")})
	m = next.(model)
	if m.filter != "ec2" || m.cursor != 0 {
		t.Fatalf("filter=%q cursor=%d", m.filter, m.cursor)
	}
	view := m.View()
	if !strings.Contains(view, "dev") || strings.Contains(view, "prod") {
		t.Fatal(view)
	}
	next, _ = m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	m = next.(model)
	if m.filter != "ec" {
		t.Fatalf("filter=%q", m.filter)
	}
	view = m.View()
	if !strings.Contains(view, "> dev") {
		t.Fatal(view)
	}
}

func TestViewNoMatch(t *testing.T) {
	m := model{hosts: sampleHosts(), filter: "zzz"}
	if !strings.Contains(m.View(), "일치하는 Host 없음") {
		t.Fatal(m.View())
	}
}

func TestInit(t *testing.T) {
	m := model{hosts: sampleHosts()}
	if cmd := m.Init(); cmd != nil {
		t.Fatal("Init should return nil")
	}
}

func TestUpdateUpDown(t *testing.T) {
	m := model{hosts: sampleHosts()}
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = next.(model)
	if m.cursor != 1 {
		t.Fatalf("cursor=%d", m.cursor)
	}
	next, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = next.(model)
	if m.cursor != 1 {
		t.Fatalf("should stay at last: %d", m.cursor)
	}
	next, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = next.(model)
	if m.cursor != 0 {
		t.Fatalf("cursor=%d", m.cursor)
	}
	next, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = next.(model)
	if m.cursor != 0 {
		t.Fatalf("should stay at first: %d", m.cursor)
	}
}

func TestViewNonDefaultPort(t *testing.T) {
	m := model{hosts: sampleHosts(), cursor: 1}
	view := m.View()
	if !strings.Contains(view, "ec2@10.0.0.2:2222") {
		t.Fatal(view)
	}
	if !strings.Contains(view, "> dev") {
		t.Fatal(view)
	}
}

func TestQQuitsWhenFilterEmpty(t *testing.T) {
	m := model{hosts: sampleHosts()}
	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	nm := next.(model)
	if !nm.quit || cmd == nil {
		t.Fatal("q should quit")
	}
}

package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"sssh/internal/sshcfg"
)

type model struct {
	hosts  []sshcfg.Host
	filter string
	cursor int
	choice *sshcfg.Host
	quit   bool
}

func Run(hosts []sshcfg.Host) (*sshcfg.Host, error) {
	if len(hosts) == 0 {
		return nil, fmt.Errorf("no Host entries in ssh config")
	}
	m := model{hosts: hosts}
	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return nil, err
	}
	out := final.(model)
	if out.quit || out.choice == nil {
		return nil, nil
	}
	return out.choice, nil
}

func (m model) Init() tea.Cmd { return nil }

func (m model) filtered() []sshcfg.Host {
	q := strings.ToLower(strings.TrimSpace(m.filter))
	if q == "" {
		return m.hosts
	}
	var out []sshcfg.Host
	for _, h := range m.hosts {
		blob := strings.ToLower(h.Alias + " " + h.HostName + " " + h.User)
		if strings.Contains(blob, q) {
			out = append(out, h)
		}
	}
	return out
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	list := m.filtered()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quit = true
			return m, tea.Quit
		case tea.KeyEnter:
			if len(list) == 0 {
				return m, nil
			}
			if m.cursor >= len(list) {
				m.cursor = len(list) - 1
			}
			h := list[m.cursor]
			m.choice = &h
			return m, tea.Quit
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(list)-1 {
				m.cursor++
			}
		case tea.KeyBackspace, tea.KeyDelete:
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.cursor = 0
			}
		case tea.KeyRunes:
			if msg.String() == "q" && m.filter == "" {
				m.quit = true
				return m, tea.Quit
			}
			m.filter += string(msg.Runes)
			m.cursor = 0
		}
	}
	return m, nil
}

func (m model) View() string {
	list := m.filtered()
	var b strings.Builder
	b.WriteString("sssh — SSH 호스트 선택 (Enter 접속, q/Esc 종료)\n")
	b.WriteString("필터: " + m.filter + "\n\n")
	if len(list) == 0 {
		b.WriteString("  (일치하는 Host 없음)\n")
		return b.String()
	}
	for i, h := range list {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		extra := h.HostName
		if h.User != "" {
			extra = h.User + "@" + extra
		}
		if h.Port != "" && h.Port != "22" {
			extra += ":" + h.Port
		}
		fmt.Fprintf(&b, "%s%s  %s\n", cursor, h.Alias, extra)
	}
	return b.String()
}

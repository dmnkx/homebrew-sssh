package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"sssh/internal/connect"
	"sssh/internal/debuglog"
	"sssh/internal/sshcfg"
	"sssh/internal/tui"
)

// options holds flags shared by the root command and subcommands.
type options struct {
	configPath string
	printCmd   bool
}

func NewRoot() *cobra.Command {
	opts := &options{}
	root := &cobra.Command{
		Use:   "sssh [alias]",
		Short: "별명으로 ~/.ssh/config Host에 SSH 접속",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := load(opts.configPath)
			if err != nil {
				return err
			}
			if len(args) == 1 {
				return connectAlias(f, args[0], opts.printCmd)
			}
			hosts := f.SelectableHosts()
			debuglog.Log("H3", "cli/root.go:root", "tui start", map[string]any{"hostCount": len(hosts), "config": f.Path})
			h, err := tui.Run(hosts)
			if err != nil {
				return err
			}
			if h == nil {
				return nil
			}
			return connectAlias(f, h.Alias, opts.printCmd)
		},
	}
	root.PersistentFlags().StringVar(&opts.configPath, "config", "", "ssh config 경로 (기본: ~/.ssh/config)")
	root.PersistentFlags().BoolVar(&opts.printCmd, "print-cmd", false, "실제 접속 대신 실행할 ssh 명령만 출력")
	root.AddCommand(listCmd(opts), addCmd(opts), editCmd(opts), rmCmd(opts), showCmd(opts), connectCmd(opts))
	return root
}

func load(configPath string) (*sshcfg.File, error) {
	p := configPath
	if p == "" {
		var err error
		p, err = sshcfg.DefaultPath()
		if err != nil {
			return nil, err
		}
	}
	return sshcfg.Load(p)
}

func connectAlias(f *sshcfg.File, alias string, printOnly bool) error {
	h, ok := f.Find(alias)
	debuglog.Log("H4", "cli/root.go:connectAlias", "lookup", map[string]any{"alias": alias, "ok": ok, "hostName": h.HostName, "user": h.User})
	if !ok {
		return fmt.Errorf("unknown host: %s", alias)
	}
	if printOnly {
		s, err := connect.PrintCmd(alias)
		if err != nil {
			return err
		}
		fmt.Println(s)
		return nil
	}
	return connect.Exec(alias)
}

func empty(s, d string) string {
	if strings.TrimSpace(s) == "" {
		return d
	}
	return s
}

package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"sssh/internal/sshcfg"
)

func addCmd(opts *options) *cobra.Command {
	var host, user, key, port, jump string
	c := &cobra.Command{
		Use:   "add <alias>",
		Short: "Host 블록 추가",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := load(opts.configPath)
			if err != nil {
				return err
			}
			alias := args[0]
			if _, exists := f.Find(alias); exists {
				return fmt.Errorf("host already exists: %s (use sssh edit)", alias)
			}
			h := sshcfg.Host{
				Alias:        alias,
				HostName:     host,
				User:         user,
				Port:         port,
				IdentityFile: key,
				ProxyJump:    jump,
			}
			if err := f.Upsert(h); err != nil {
				return err
			}
			return f.Save()
		},
	}
	c.Flags().StringVar(&host, "host", "", "HostName (IP 또는 DNS)")
	c.Flags().StringVar(&user, "user", "", "User")
	c.Flags().StringVar(&key, "key", "", "IdentityFile")
	c.Flags().StringVar(&port, "port", "", "Port")
	c.Flags().StringVar(&jump, "jump", "", "ProxyJump")
	_ = c.MarkFlagRequired("host")
	return c
}

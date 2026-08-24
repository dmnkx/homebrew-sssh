package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func editCmd(opts *options) *cobra.Command {
	var host, user, key, port, jump string
	c := &cobra.Command{
		Use:   "edit <alias>",
		Short: "Host 필드 갱신",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := load(opts.configPath)
			if err != nil {
				return err
			}
			h, ok := f.Find(args[0])
			if !ok {
				return fmt.Errorf("unknown host: %s", args[0])
			}
			if cmd.Flags().Changed("host") {
				h.HostName = host
			}
			if cmd.Flags().Changed("user") {
				h.User = user
			}
			if cmd.Flags().Changed("key") {
				h.IdentityFile = key
			}
			if cmd.Flags().Changed("port") {
				h.Port = port
			}
			if cmd.Flags().Changed("jump") {
				h.ProxyJump = jump
			}
			if err := f.Upsert(h); err != nil {
				return err
			}
			return f.Save()
		},
	}
	c.Flags().StringVar(&host, "host", "", "HostName")
	c.Flags().StringVar(&user, "user", "", "User")
	c.Flags().StringVar(&key, "key", "", "IdentityFile")
	c.Flags().StringVar(&port, "port", "", "Port")
	c.Flags().StringVar(&jump, "jump", "", "ProxyJump")
	return c
}

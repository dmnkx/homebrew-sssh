package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func listCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Host 목록",
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := load(opts.configPath)
			if err != nil {
				return err
			}
			for _, h := range f.SelectableHosts() {
				port := h.Port
				if port == "" {
					port = "22"
				}
				fmt.Printf("%-20s %s@%s:%s %s\n", h.Alias, empty(h.User, "-"), empty(h.HostName, "-"), port, h.IdentityFile)
			}
			return nil
		},
	}
}

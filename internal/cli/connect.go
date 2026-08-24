package cli

import (
	"github.com/spf13/cobra"
)

func connectCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "connect <alias>",
		Short: "별명으로 ssh 접속",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := load(opts.configPath)
			if err != nil {
				return err
			}
			return connectAlias(f, args[0], opts.printCmd)
		},
	}
}

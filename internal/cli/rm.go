package cli

import (
	"github.com/spf13/cobra"
)

func rmCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "rm <alias>",
		Short: "Host 블록 제거",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := load(opts.configPath)
			if err != nil {
				return err
			}
			if err := f.Remove(args[0]); err != nil {
				return err
			}
			return f.Save()
		},
	}
}

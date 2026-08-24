package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func showCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "show <alias>",
		Short: "Host 블록 출력",
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
			fmt.Println(h.BlockString())
			return nil
		},
	}
}

package cmd

import (
	"github.com/libmojito/mingle/convo"
	"github.com/spf13/cobra"
)

func NewThreadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "thread",
		Short: "Create a new thread",
		Long:  `Create a new thread`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("%s\n", convo.NewThread().ThreadID)
		},
	}

	return cmd
}

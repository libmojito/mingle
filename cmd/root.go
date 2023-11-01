package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mingle",
		Short: "Talk about stuff",
		Long:  `Talk about stuff`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(NewChatCmd())
	cmd.AddCommand(NewThreadCmd())
	cmd.AddCommand(NewServeCmd())

	return cmd
}

func Execute() {
	err := NewCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}

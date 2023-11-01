package cmd

import (
	"github.com/libmojito/mingle/convo"
	"github.com/spf13/cobra"
)

func NewServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "a simple chat service",
		Long:  `a simple chat service`,
		Run: func(cmd *cobra.Command, args []string) {
			svr := convo.NewServer(convo.WithAddress(":8080"))
			svr.Run()
		},
	}

	return cmd
}

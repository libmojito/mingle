package cmd

import (
	"log"

	"github.com/libmojito/mingle/convo"
	"github.com/spf13/cobra"
)

func NewThreadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "thread",
		Short: "Create a new thread",
		Long:  `Create a new thread`,
		Run: func(cmd *cobra.Command, args []string) {
			svrAddr, err := cmd.Flags().GetString(FlagServerAddress)
			if err != nil {
				log.Fatal(err)
			}
			cmd.Printf(
				"%s\n",
				convo.NewThreadClient(
					convo.WithServerAddress(svrAddr),
				).NewThread(),
			)
		},
	}

	cmd.Flags().String(FlagServerAddress, convo.DefaultServerAddress, "The server address")

	return cmd
}

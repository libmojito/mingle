package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/libmojito/mingle/convo"
	"github.com/spf13/cobra"
)

const (
	FlagThread      = "thread"
	FlagContent     = "content"
	FlagInteractive = "interactive"
)

func NewChatCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chat",
		Short: "Talk about stuff",
		Long:  `Talk about stuff`,
		Run: func(cmd *cobra.Command, args []string) {

			threadID, err := cmd.Flags().GetString(FlagThread)
			if err != nil {
				log.Fatal(err)
			}
			content, err := cmd.Flags().GetString(FlagContent)
			if err != nil {
				log.Fatal(err)
			}
			interactive, err := cmd.Flags().GetBool(FlagInteractive)
			if err != nil {
				log.Fatal(err)
			}

			clt := convo.NewClient(threadID)

			if interactive {
				if len(os.Getenv("DEBUG")) > 0 {
					f, err := tea.LogToFile("debug.log", "debug")
					if err != nil {
						fmt.Println("fatal:", err)
						os.Exit(1)
					}
					defer f.Close()
				}

				p := tea.NewProgram(convo.InitialModel(clt), tea.WithAltScreen())

				if _, err := p.Run(); err != nil {
					log.Fatal(err)
				}

				return
			}

			msg := clt.SendMessageContent(content).Content

			cmd.Println("response: " + msg)

			// message := strings.Join(args, " ")
			cmd.Printf("Hi, your thread ID is %s\n", threadID)
		},
	}

	cmd.Flags().String(
		FlagThread,
		"",
		"the required thread ID",
	)
	if err := cmd.MarkFlagRequired(FlagThread); err != nil {
		log.Fatal(err)
	}

	cmd.Flags().Bool(FlagInteractive, false, "interactive session")

	cmd.Flags().String(FlagContent, "", "the message content")

	return cmd
}

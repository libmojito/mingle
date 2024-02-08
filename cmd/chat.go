package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/libmojito/mingle/convo"
	"github.com/spf13/cobra"
)

const (
	FlagThread      = "thread"
	FlagContent     = "content"
	FlagSecurity    = "security"
	FlagInteractive = "interactive"
)

func parseSecurity(security []string) map[string]string {
	tokens := make(map[string]string)
	for _, s := range security {
		ts := strings.Split(s, "=")
		if len(ts) == 1 {
			tokens["__default__"] = s
		} else {
			k := ts[0]
			v := strings.Join(ts[1:], "=")
			tokens[k] = v
		}
	}
	return tokens
}

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
			security, err := cmd.Flags().GetStringArray(FlagSecurity)
			if err != nil {
				log.Fatal(err)
			}
			tokens := parseSecurity(security)

			interactive, err := cmd.Flags().GetBool(FlagInteractive)
			if err != nil {
				log.Fatal(err)
			}

			clt := convo.NewClient(threadID, tokens)

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
	cmd.Flags().StringArray(FlagSecurity, []string{}, "the bearer tokens to be added for different ")

	return cmd
}

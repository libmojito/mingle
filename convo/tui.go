package convo

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

const (
	viewWidth      = 80
	textareaHeight = 2
)

type errMsg error

type model struct {
	viewport       viewport.Model
	viewportRows   int
	messages       []string // rendered messages
	textarea       textarea.Model
	textareaRows   int
	senderStyle    lipgloss.Style
	responderStyle lipgloss.Style
	client         *Client
	spinner        spinner.Model
	blur           bool
	err            error
}

func InitialModel(clt *Client) model {
	ta := textarea.New()
	ta.Placeholder = "Ctrl-S or Shift-Right to send a message..."
	ta.Focus()

	ta.Prompt = "| "
	ta.CharLimit = 500

	ta.SetWidth(viewWidth)
	ta.SetHeight(textareaHeight)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	vp := viewport.New(viewWidth, 25)
	vp.SetContent(`Welcom to the mingle!
Type a message and press Enter to send.`)

	// no need to send the last newline in the message
	ta.KeyMap.InsertNewline.SetEnabled(true)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		textarea:       ta,
		textareaRows:   textareaHeight,
		messages:       []string{},
		viewport:       vp,
		viewportRows:   25,
		senderStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		responderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("12")),
		spinner:        s,
		blur:           false,
		client:         clt,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc, tea.KeyCtrlD:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyCtrlS, tea.KeyShiftRight:
			c := m.textarea.Value()
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+c)

			reply := m.getResponse(c)

			str := wordwrap.String(strings.Join(m.messages, "\n"), viewWidth)
			m.viewport.SetContent(str)
			m.textarea.Reset()
			m.viewport.GotoBottom()

			// adjust the size of viewport and textarea
			offset := m.textarea.Height() - textareaHeight
			m.textarea.SetHeight(textareaHeight)
			m.textareaRows = m.textarea.Height()
			m.viewport.Height = m.viewport.Height + offset
			m.viewportRows = m.viewport.Height

			m.textarea.Blur()
			m.blur = true

			return m, tea.Batch(m.spinner.Tick, reply)
		case tea.KeyEnter:
			m.textarea.SetHeight(m.textarea.Height() + 1)
			m.textareaRows = m.textarea.Height()
			m.viewport.Height = m.viewport.Height - 1
			m.viewportRows = m.viewport.Height
		default:
			m.textarea.SetHeight(m.textareaRows + m.textarea.LineInfo().RowOffset)
			m.viewport.Height = m.viewportRows - m.textarea.LineInfo().RowOffset
		}
	case spinner.TickMsg:

		if m.blur {
			var cmd tea.Cmd

			m.spinner, cmd = m.spinner.Update(msg)
			str := wordwrap.String(
				strings.Join(m.messages, "\n")+
					"\n"+
					m.responderStyle.Render(m.spinner.View())+
					"Waiting for response...",
				viewWidth,
			)
			m.viewport.SetContent(str)

			return m, tea.Batch(tiCmd, vpCmd, cmd)
		}

	case Message:
		m.messages = append(m.messages, m.responderStyle.Render("Assistant: ")+msg.Content)
		str := wordwrap.String(
			strings.Join(m.messages, "\n"),
			viewWidth,
		)
		m.textarea.Focus()
		m.blur = false
		m.viewport.SetContent(str)

	case errMsg:
		m.err = msg
		return m, nil

	case tea.WindowSizeMsg:
		m.viewport.Height = msg.Height - m.textarea.Height() - 3
		m.viewportRows = m.viewport.Height + m.textarea.LineInfo().RowOffset
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) getResponse(chat string) tea.Cmd {
	return func() tea.Msg {
		return m.client.SendMessageContent(chat)
	}
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

package ui

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/chrispeterkins/claude-history/internal/data"
)

// exportMarkdown renders the current conversation as clean markdown.
func exportMarkdown(messages []data.Message) string {
	var sb strings.Builder

	for _, msg := range messages {
		switch msg.Type {
		case "user":
			if msg.RawText == "" {
				continue
			}
			sb.WriteString("## You\n")
			sb.WriteString(fmt.Sprintf("*%s*\n\n", msg.Timestamp.Format("Jan 02 15:04:05")))
			sb.WriteString(msg.RawText)
			sb.WriteString("\n\n---\n\n")

		case "assistant":
			sb.WriteString("## Claude\n")
			sb.WriteString(fmt.Sprintf("*%s*", msg.Timestamp.Format("Jan 02 15:04:05")))
			if msg.Model != "" {
				sb.WriteString(fmt.Sprintf(" · *%s*", msg.Model))
			}
			sb.WriteString("\n\n")

			// Tool calls
			for _, pair := range msg.ToolPairs {
				sb.WriteString(fmt.Sprintf("**[%s]**", pair.Name))
				if pair.Use.Input != nil {
					if cmd, ok := pair.Use.Input["command"].(string); ok {
						sb.WriteString(fmt.Sprintf(" `%s`", cmd))
					} else if path, ok := pair.Use.Input["file_path"].(string); ok {
						sb.WriteString(fmt.Sprintf(" `%s`", path))
					}
				}
				sb.WriteString("\n")
				if pair.Result.Content != "" {
					content := pair.Result.Content
					if len(content) > 1000 {
						content = content[:1000] + "\n... (truncated)"
					}
					sb.WriteString("```\n")
					sb.WriteString(content)
					sb.WriteString("\n```\n")
				}
				sb.WriteString("\n")
			}

			// Text content
			sb.WriteString(msg.RawText)
			sb.WriteString("\n\n---\n\n")

		case "system":
			if msg.Subtype == "turn_duration" && msg.DurationMs > 0 {
				dur := time.Duration(msg.DurationMs) * time.Millisecond
				sb.WriteString(fmt.Sprintf("*Turn: %s*\n\n", dur.Round(time.Millisecond)))
			}
		}
	}

	return sb.String()
}

// copyToClipboard copies text to the system clipboard using pbcopy (macOS).
func copyToClipboard(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

type clipboardCopiedMsg struct {
	err error
}

func (m Model) copyConversationCmd() tea.Cmd {
	messages := m.messages
	return func() tea.Msg {
		md := exportMarkdown(messages)
		err := copyToClipboard(md)
		return clipboardCopiedMsg{err: err}
	}
}

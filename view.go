package deployaroo

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func boxString(s string, width int) string {
	output := "\x1b[2m│\x1b[0m "
	if len(s) > width {
		output += s[0:width - 1]
		output += "…"
	} else {
		output += s
		for i := 0; i < width - len(s); i++ {
			output += " "
		}
	}
	output += " \x1b[2m│\x1b[0m"
	return output
}

func (m model) renderDeployments() string {
	output := ""
	renderedDeployments := [][]string{}
    for _, deployment := range m.deployments {
		rule := ""
		for i := 0; i < m.boxWidth + 2; i++ {
			rule += "─"
		}
		box := []string{
			"\x1b[2m╭" + rule + "╮\x1b[0m",
			boxString(deployment.name, m.boxWidth),
			boxString(deployment.ticket, m.boxWidth),
			boxString(deployment.version, m.boxWidth),
			boxString(deployment.url, m.boxWidth),
			boxString("\x1b[32m♻︎\x1b[0m" + deployment.lastRestart.String(), m.boxWidth),
			boxString("\x1b[34m⏲\x1b[0m", m.boxWidth),
			"\x1b[2m╰" + rule + "╯\x1b[0m",
		}
		renderedDeployments = append(renderedDeployments, box)
    }
	if len(renderedDeployments) > 0 {
		for i := 0; i < len(renderedDeployments[0]); i++ {
			for _, rd := range renderedDeployments {
				output += rd[i]
				output += "   "
			}
			output += "\n"
		}
	}
	return output
}

func (m model) defaultView() string {
	s := ""

	// Render the list of deployments.
	s += m.renderDeployments();

	// Render the history.
	for _, message := range m.history[:3] {
		s += fmt.Sprintf("\n    \x1b[2m%s\x1b[0m\n", message)
	}

	// Render the text input.
	s += "\n\x1b[2m╭"
	for i := 0; i < m.termWidth - 2; i++ {
		s += "─"
	}
	s += "╮\n│ \x1b[0m"
	s += m.textInput.View()
	s += "\x1b[2m │\n╰"
	for i := 0; i < m.termWidth - 2; i++ {
		s += "─"
	}
	s += "╯\x1b[0m"

	s += m.spinner.View()

	return s
}

func (m model) errorView() string {
	s := "Error:\n\n"
	return s
}

func (m model) View() string {
    // Render the header.
    s := "\nDeployaroo!\n\n"

    s += m.defaultView()

    // Render the footer.
    s += "\n\n    \x1b[2mPress (esc) to quit.\x1b[0m\n\n"

    return s
}

func (m model) updateSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.termHeight = msg.Height
	m.termWidth = msg.Width
	m.textInput.Width = msg.Width - 7
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var dCmd tea.Cmd
    var sCmd tea.Cmd
    var tiCmd tea.Cmd

	if len(m.deployments) <= 0 && !m.loadingDeployments {
		dCmd = m.FetchDeployments()
	}
	
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		return m.updateSize(msg)

    case tea.KeyMsg:

        switch msg.Type {

        case tea.KeyCtrlC, tea.KeyEsc:
            return m, tea.Quit

        case tea.KeyEnter:
			message := m.textInput.Value()
			if len(message) > 0 {
				m.history = append(m.history, message)
				m.history = m.history[1:]
			}
        }
    }

	m.textInput, tiCmd = m.textInput.Update(msg)
	m.spinner, sCmd = m.spinner.Update(msg)
    return m, tea.Batch(dCmd, sCmd, tiCmd)
}

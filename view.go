package deployaroo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Leo-li-dotmatics/dmngr"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mitchellh/go-wordwrap"
	"github.com/paul-freeman/deployerai"
)

func (m model) View() string {
	b := strings.Builder{}

	// Title
	b.WriteString(TitleStyle.Width(m.termWidth).Render("Deployaroo! | ", string(m.view)))

	// Content
	switch m.view {
	case LoadingView:
		m.loadingView(&b)
	case InputView:
		m.inputView(&b)
	case ChoiceView:
		m.choiceView(&b)
	case ErrorView:
		m.errorView(&b)
	default:
		panic(fmt.Sprintf("unknown view: %v", m.view))
	}

	// Footer
	b.WriteString(TitleStyle.Width(m.termWidth).Render("Press (esc) to quit."))

	return b.String()
}

func (m model) inputView(b *strings.Builder) {
	b.WriteString("\n\n")
	b.WriteString(m.textInput.View())
}

func (m model) errorView(b *strings.Builder) {
	b.WriteString(SubtitleStyle.Width(m.termWidth).Padding(m.termHeight/2-4, 4).Render("Error: ", m.err.Error()))
}

func (m model) choiceView(b *strings.Builder) {
	wrapped := wordwrap.WrapString(m.choice.Message, uint(m.termWidth-10))
	b.WriteString(SubtitleStyle.Width(m.termWidth).Height(m.termHeight-6).Align(lipgloss.Left).Padding(6, 4).Render(wrapped))
}

func (m model) loadingView(b *strings.Builder) {
	b.WriteString(LoadingStyle.Width(m.termWidth).Padding(m.termHeight/2-4, 4).Render("Finding best deployment target", m.spinner.View()))
}

func (m model) updateSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.termHeight = msg.Height
	m.termWidth = msg.Width
	m.textInput.SetWidth(msg.Width)
	m.textInput.SetHeight(msg.Height - 6)
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var tiCmd tea.Cmd
	m.textInput, tiCmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.updateSize(msg)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.userInput = m.textInput.Value()
			m.view = LoadingView
			return m, tea.Batch(m.spinner.Tick, dmngr.GetAllClustersInfo(context.TODO()))
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case dmngr.AllClustersInfoResp:
		m.deployments = msg.Targets
		m.deploymentTable = m.buildTable(msg.Targets)
		m.view = InputView
		req := deployerai.Request{
			MessageFromUser: m.userInput,
			AdditionalNotes: "",
		}
		for _, d := range m.deployments {
			req.DeploymentTargets = append(req.DeploymentTargets, deployerai.Target{
				Name:                       d.Cluster,
				CurrentImage:               d.CurrentImage,
				CurrentImageDeploymentTime: d.LastImageUpdate,
				LastRestart:                d.LastRestart,
				LastUsed:                   d.LastLogTime,
			})
		}
		m.view = LoadingView
		return m, deployerai.ChooseDeploymentTarget(context.TODO(), deployerai.ModelGPT4, req)
	case dmngr.Error:
		m.err = errors.New(msg.Message)
		m.view = ErrorView
		return m, nil
	case deployerai.Choice:
		m.choice = msg
		m.view = ChoiceView
	case deployerai.Error:
		m.err = errors.New(msg.Message)
		m.view = ErrorView
	default:
		// panic(fmt.Sprintf("unknown message: %T", msg))
	}

	return m, tiCmd
}

func (m model) buildTable(deployments []dmngr.Target) table.Model {
	columns := []table.Column{
		{Title: "Name", Width: 30},
		{Title: "Cluster", Width: 30},
		{Title: "Image", Width: 30},
		{Title: "Last Log", Width: 20},
		{Title: "Last Update", Width: 20},
		{Title: "Last Restart", Width: 20},
	}

	rows := make([]table.Row, len(deployments))
	for i, d := range deployments {
		rows[i] = table.Row{
			d.Name,
			d.Cluster,
			d.CurrentImage,
			d.LastLogTime.String(),
			d.LastImageUpdate.String(),
			d.LastRestart.String(),
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithWidth(m.termWidth-2),
		table.WithHeight(m.termHeight-8),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}

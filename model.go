package deployaroo

import (
	"time"

	"github.com/Leo-li-dotmatics/dmngr"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/paul-freeman/deployerai"
)

type deployment struct {
	name        string
	version     string
	ticket      string
	url         string
	pr          string
	lastRestart time.Time
	lastLog     time.Time
	lastUpdate  time.Time
}

type model struct {
	deployments        []dmngr.Target
	deploymentTable    table.Model
	history            []string
	textInput          textarea.Model
	spinner            spinner.Model
	termHeight         int
	termWidth          int
	boxWidth           int
	loadingDeployments bool
	loadingResponse    bool
	view               view
	userInput          string
	choice             deployerai.Choice

	err error
}

type view string

const (
	LoadingView view = "LoadingView"
	InputView        = "InputView"
	ChoiceView       = "ChoiceView"
	ErrorView        = "ErrorView"
)

func InitialModel() model {
	textInput := textarea.New()
	textInput.Placeholder = "Tell me what you want to deploy!"
	textInput.Focus()
	textInput.Prompt = "┃ "
	textInput.CharLimit = 280
	textInput.FocusedStyle.CursorLine = lipgloss.NewStyle()
	textInput.ShowLineNumbers = false
	textInput.KeyMap.InsertNewline.SetEnabled(false)

	boxWidth := 30
	m := model{
		deploymentTable:    table.New(),
		history:            []string{"", "", "What would you like to do?"},
		textInput:          textInput,
		spinner:            spinner.New(),
		boxWidth:           boxWidth,
		loadingDeployments: false,
		loadingResponse:    false,
		view:               InputView,
		err:                nil,
	}
	m.spinner.Spinner = spinner.Moon
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink)
}

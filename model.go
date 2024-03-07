package deployaroo

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type deployment struct {
	name string
	version string
	ticket string
	url string
	pr string
	lastRestart time.Time
	lastLog time.Time
	lastUpdate time.Time
}

type model struct {
	deployments []deployment
	history []string
	textInput textinput.Model
	spinner spinner.Model
	termHeight int
	termWidth  int
	boxWidth int
	loadingDeployments bool
	loadingResponse bool
	view int
}

const (
	DefaultView = iota
	ErrorView
)

func InitialModel() model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 0
	ti.Width = 70

	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lg.NewStyle().Foreground(lg.Color("205"))

	return model {
		deployments: []deployment{},
		history: []string{"", "", "What would you like to do?"},
		textInput: ti,
		spinner: s,
		boxWidth: 30,
		loadingDeployments: false,
		loadingResponse: false,
		view: DefaultView,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

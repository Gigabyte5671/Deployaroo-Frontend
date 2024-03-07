package deployaroo

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
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


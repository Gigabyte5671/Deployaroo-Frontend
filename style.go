package deployaroo

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultWidth = 75
)

var ListStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FAFAFA")).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#2E8B57")).
	Padding(0, 1).
	Margin(1)

var TitleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF00FF")).
	Background(lipgloss.Color("#39FF14")).
	Bold(true).
	Padding(1, 1).
	Width(defaultWidth).
	Align(lipgloss.Center)

var FooterStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF00FF")).
	Background(lipgloss.Color("#39FF14")).
	Bold(true).
	Width(defaultWidth).
	Align(lipgloss.Center)

var TextInputStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF4500")).
	Width(defaultWidth - 4).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#2E8B57")).
	Padding(1).
	Margin(1)

var SubtitleStyle = lipgloss.NewStyle().
	Padding(1, 1).
	Width(defaultWidth).
	Align(lipgloss.Center)

var LoadingStyle = lipgloss.NewStyle().
	Padding(1, 1).
	Width(defaultWidth).
	Align(lipgloss.Center)

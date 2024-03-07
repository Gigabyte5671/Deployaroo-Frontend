package main

import (
	"fmt"
	"os"

	deployaroo "github.com/Gigabyte5671/Deployaroo-Frontend"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
    p := tea.NewProgram(deployaroo.InitialModel(), tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}

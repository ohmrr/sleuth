package cmd

import "charm.land/lipgloss/v2"

var successStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#3CD649"))

var errorStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FC6749"))

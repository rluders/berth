// Package tui provides the Terminal User Interface for Berth.
package tui

import "github.com/charmbracelet/lipgloss"

// Theme defines the color and styling for the application.
type Theme struct {
	AppStyle           lipgloss.Style
	HeaderStyle        lipgloss.Style
	FooterStyle        lipgloss.Style
	StatusMessageStyle lipgloss.Style
	TableSelectedStyle lipgloss.Style
	TableHeaderStyle   lipgloss.Style
}

// DefaultTheme returns a new Theme with default styles.
func DefaultTheme() Theme {
	return Theme{
		AppStyle:           lipgloss.NewStyle().Padding(1, 2),
		HeaderStyle:        lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(0, 1).Border(lipgloss.NormalBorder(), false, false, true, false),
		FooterStyle:        lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(0, 1).Border(lipgloss.NormalBorder(), true, false, false, false),
		StatusMessageStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Padding(0, 1),
		TableSelectedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false),
		TableHeaderStyle:   lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(false),
	}
}
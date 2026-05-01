package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Catppuccin Mocha palette
const (
	colorBase    = "#1e1e2e"
	colorMantle  = "#181825"
	colorSurface = "#313244"
	colorOverlay = "#45475a"
	colorText    = "#cdd6f4"
	colorSubtext = "#a6adc8"
	colorMuted   = "#6c7086"

	colorMauve  = "#cba6f7"
	colorBlue   = "#89b4fa"
	colorSky    = "#89dceb"
	colorGreen  = "#a6e3a1"
	colorYellow = "#f9e2af"
	colorPeach  = "#fab387"
	colorRed    = "#f38ba8"
	colorTeal   = "#94e2d5"
	colorLavend = "#b4befe"

	colorCrust = "#11111b"
)

// Theme defines all visual styles for the application.
type Theme struct {
	// App chrome
	AppStyle  lipgloss.Style
	AppBg     lipgloss.Style

	// Header
	HeaderStyle     lipgloss.Style
	HeaderLogoStyle lipgloss.Style
	HeaderEngStyle  lipgloss.Style

	// Tabs
	TabBarStyle    lipgloss.Style
	ActiveTabStyle lipgloss.Style
	InactiveTabStyle lipgloss.Style
	TabCountStyle  lipgloss.Style

	// Footer
	FooterStyle      lipgloss.Style
	FooterKeyStyle   lipgloss.Style
	FooterDescStyle  lipgloss.Style

	// Status
	StatusMessageStyle lipgloss.Style
	StatusOKStyle      lipgloss.Style
	StatusErrStyle     lipgloss.Style
	SpinnerStyle       lipgloss.Style

	// Tables
	TableHeaderStyle   lipgloss.Style
	TableSelectedStyle lipgloss.Style
	TableRowStyle      lipgloss.Style
	TableRowAltStyle   lipgloss.Style

	// Badges
	BadgeRunningStyle   lipgloss.Style
	BadgeStoppedStyle   lipgloss.Style
	BadgePausedStyle    lipgloss.Style
	BadgeRestartStyle   lipgloss.Style
	BadgeCreatedStyle   lipgloss.Style

	// Cards (details view)
	CardStyle      lipgloss.Style
	CardTitleStyle lipgloss.Style
	CardValueStyle lipgloss.Style
	SectionStyle   lipgloss.Style

	// Modal
	ModalOverlayStyle lipgloss.Style
	ModalBoxStyle     lipgloss.Style
	ModalTitleStyle   lipgloss.Style
	ModalBodyStyle    lipgloss.Style

	// Buttons
	ButtonPrimaryStyle   lipgloss.Style
	ButtonDangerStyle    lipgloss.Style
	ButtonSecondaryStyle lipgloss.Style
	ButtonFocusedStyle   lipgloss.Style

	// Filter input
	FilterStyle lipgloss.Style

	// Log viewer
	LogTimestampStyle lipgloss.Style
	LogErrorStyle     lipgloss.Style
	LogWarnStyle      lipgloss.Style
	LogInfoStyle      lipgloss.Style
	LogDebugStyle     lipgloss.Style
	LogLineNumStyle   lipgloss.Style
	LogFollowStyle    lipgloss.Style
	LogPausedStyle    lipgloss.Style

	// Viewport
	ViewportStyle lipgloss.Style

	// Dividers
	DividerStyle lipgloss.Style

	// Container accordion
	GroupHeaderStyle lipgloss.Style
	GroupChildStyle  lipgloss.Style

	// Legacy (referenced by view.go / update.go)
	ModalStyle lipgloss.Style
}

// DefaultTheme returns a new Theme using Catppuccin Mocha palette.
func DefaultTheme() Theme {
	t := Theme{}

	// App chrome
	t.AppStyle = lipgloss.NewStyle().Padding(0, 0)
	t.AppBg = lipgloss.NewStyle().Background(lipgloss.Color(colorBase))

	// Header
	t.HeaderStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(colorMantle)).
		Foreground(lipgloss.Color(colorText)).
		Padding(0, 2).
		Bold(false)
	t.HeaderLogoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMauve)).
		Bold(true)
	t.HeaderEngStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorSubtext)).
		Background(lipgloss.Color(colorSurface)).
		Padding(0, 1)

	// Tabs
	t.TabBarStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(colorCrust)).
		Padding(0, 0)
	t.ActiveTabStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMauve)).
		Background(lipgloss.Color(colorBase)).
		Bold(true).
		Padding(0, 2).
		Border(lipgloss.Border{Bottom: "▔"}, false, false, true, false).
		BorderForeground(lipgloss.Color(colorMauve))
	t.InactiveTabStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMuted)).
		Background(lipgloss.Color(colorCrust)).
		Padding(0, 2)
	t.TabCountStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorOverlay)).
		Background(lipgloss.Color(colorSurface)).
		Padding(0, 1)

	// Footer
	t.FooterStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(colorMantle)).
		Foreground(lipgloss.Color(colorMuted)).
		Padding(0, 2)
	t.FooterKeyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMauve)).
		Bold(true)
	t.FooterDescStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorSubtext))

	// Status
	t.StatusMessageStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorYellow)).
		Padding(0, 2)
	t.StatusOKStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGreen)).
		Padding(0, 2)
	t.StatusErrStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorRed)).
		Padding(0, 2)
	t.SpinnerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMauve))

	// Tables
	t.TableHeaderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMauve)).
		Background(lipgloss.Color(colorMantle)).
		Bold(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(colorSurface)).
		BorderBottom(true)
	t.TableSelectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBase)).
		Background(lipgloss.Color(colorMauve)).
		Bold(true)
	t.TableRowStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText))
	t.TableRowAltStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText)).
		Background(lipgloss.Color(colorMantle))

	// Badges
	t.BadgeRunningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBase)).
		Background(lipgloss.Color(colorGreen)).
		Padding(0, 1).
		Bold(true)
	t.BadgeStoppedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBase)).
		Background(lipgloss.Color(colorRed)).
		Padding(0, 1).
		Bold(true)
	t.BadgePausedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBase)).
		Background(lipgloss.Color(colorYellow)).
		Padding(0, 1).
		Bold(true)
	t.BadgeRestartStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBase)).
		Background(lipgloss.Color(colorBlue)).
		Padding(0, 1).
		Bold(true)
	t.BadgeCreatedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBase)).
		Background(lipgloss.Color(colorTeal)).
		Padding(0, 1).
		Bold(true)

	// Cards
	t.CardStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colorSurface)).
		Padding(0, 1).
		Margin(0, 0, 1, 0)
	t.CardTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMauve)).
		Bold(true)
	t.CardValueStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText))
	t.SectionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorLavend)).
		Bold(true).
		MarginTop(1)

	// Modal
	t.ModalBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colorMauve)).
		Background(lipgloss.Color(colorMantle)).
		Padding(1, 3)
	t.ModalTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMauve)).
		Bold(true).
		MarginBottom(1)
	t.ModalBodyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText))
	// Legacy alias
	t.ModalStyle = t.ModalBoxStyle

	// Buttons
	t.ButtonPrimaryStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBase)).
		Background(lipgloss.Color(colorBlue)).
		Padding(0, 2).
		Margin(0, 1)
	t.ButtonDangerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBase)).
		Background(lipgloss.Color(colorRed)).
		Padding(0, 2).
		Margin(0, 1)
	t.ButtonSecondaryStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText)).
		Background(lipgloss.Color(colorSurface)).
		Padding(0, 2).
		Margin(0, 1)
	t.ButtonFocusedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBase)).
		Background(lipgloss.Color(colorMauve)).
		Padding(0, 2).
		Margin(0, 1).
		Bold(true)

	// Filter
	t.FilterStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText)).
		Background(lipgloss.Color(colorSurface)).
		Padding(0, 1)

	// Log viewer
	t.LogTimestampStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorMuted))
	t.LogErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorRed)).Bold(true)
	t.LogWarnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorYellow))
	t.LogInfoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorBlue))
	t.LogDebugStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorMuted))
	t.LogLineNumStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorOverlay))
	t.LogFollowStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBase)).
		Background(lipgloss.Color(colorGreen)).
		Padding(0, 1).
		Bold(true)
	t.LogPausedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBase)).
		Background(lipgloss.Color(colorYellow)).
		Padding(0, 1)

	// Viewport
	t.ViewportStyle = lipgloss.NewStyle().
		Padding(0, 1)

	// Divider
	t.DividerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorSurface))

	// Container accordion
	t.GroupHeaderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMauve)).
		Bold(true)
	t.GroupChildStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorSubtext))

	return t
}

// StatusBadge returns a styled status badge string for a container status.
func StatusBadge(status string) string {
	switch {
	case strings.HasPrefix(status, "Up"), status == "running":
		return currentTheme.BadgeRunningStyle.Render("▶ " + status)
	case status == "paused":
		return currentTheme.BadgePausedStyle.Render("⏸ " + status)
	case status == "restarting":
		return currentTheme.BadgeRestartStyle.Render("↻ " + status)
	case status == "created":
		return currentTheme.BadgeCreatedStyle.Render("● " + status)
	default:
		return currentTheme.BadgeStoppedStyle.Render("■ " + status)
	}
}

// StatusColor returns a lipgloss-styled string for a container status (plain text, no badge).
func StatusColor(status string) string {
	switch {
	case strings.HasPrefix(status, "Up"), status == "running":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorGreen)).Render(status)
	case status == "paused":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorYellow)).Render(status)
	case status == "restarting":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorBlue)).Render(status)
	case status == "created":
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorTeal)).Render(status)
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colorRed)).Render(status)
	}
}

// FooterHint renders a key+description pair for the footer.
func FooterHint(k, desc string) string {
	return currentTheme.FooterKeyStyle.Render(k) +
		currentTheme.FooterDescStyle.Render(" "+desc)
}

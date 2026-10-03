// Package ui renders Moki's terminal interface.
package ui

import "github.com/charmbracelet/lipgloss"

var (
	// accent is Moki's signature colour.
	accent = lipgloss.Color("205")
	dim    = lipgloss.Color("241")
	green  = lipgloss.Color("42")

	headerStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)

	taglineStyle = lipgloss.NewStyle().
			Foreground(dim)

	userStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	mokiStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	helpStyle = lipgloss.NewStyle().
			Foreground(dim)

	statusStyle = lipgloss.NewStyle().
			Foreground(dim).
			Italic(true)
)

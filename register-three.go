package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RegisterThreeModel struct {
	page             int
	steps            int
	current          int
	err              error
	headerStyle      lipgloss.Style
	subheaderStyle   lipgloss.Style
	username         string
	email            string
	birth            string
}

func getPageThreeModel(renderer *lipgloss.Renderer, steps int) RegisterThreeModel {
	headerStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	subheaderStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
	return RegisterThreeModel {
		page: 3,
		steps: steps,
		err:       nil,
		headerStyle: headerStyle,
		subheaderStyle: subheaderStyle,
		current: 0,
	}
}

func (m RegisterThreeModel) Init() tea.Cmd {
	return nil
}

type AcceptMsg struct {}

func accept() tea.Msg {
	return AcceptMsg{}
}

func (m RegisterThreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd 

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, nil;
		case "enter": 
			return m, accept
		}
	case error:
		m.err = msg
		return m, nil
	}

	return m, cmd
}

func (m RegisterThreeModel) View() string {
	buttonStyle := lipgloss.NewStyle().
		MarginTop(2).
		Width(30).
		Align(lipgloss.Right)
	
	button := buttonStyle.
		Render(getButtonPrefix(true) + "[ Confirm ]")

	return m.headerStyle.Render("Confirm data") + "\n" +
		m.subheaderStyle.Render(fmt.Sprintf("Step %d of %d", m.page, m.steps)) +
		"\n\n" +
		"Username: " + m.username + "\n" +
		"E-mail: " + m.email + "\n" +
		"Date of birth: " + m.birth + "\n" +
		button

}

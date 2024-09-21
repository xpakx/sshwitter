package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func getRegisterModel(renderer *lipgloss.Renderer, username string) (RegisterModel) {
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
	usernameStyle := renderer.NewStyle().Foreground(lipgloss.Color("5"))
	pageStyle := renderer.NewStyle().
		Height(5).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1, 2)


	activeDot := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	inactiveDot := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")

	pages := []tea.Model{
		getPageOneModel(3, username),
		getPageTwoModel(3),
		getPageOneModel(3, ""),
	}
	
	return RegisterModel{ 
		name: "sshwitter", 
		username: username,
		txtStyle: txtStyle, 
		quitStyle: quitStyle,
		userStyle: usernameStyle,
		pageStyle: pageStyle,
		activeDot: activeDot,
		inactiveDot: inactiveDot,
		currentView: 0,
		pages: pages,
	}
}

type RegisterModel struct {
	name         string
	username     string
	txtStyle     lipgloss.Style
	quitStyle    lipgloss.Style
	userStyle    lipgloss.Style
	pageStyle    lipgloss.Style
	activeDot    string
	inactiveDot  string
	currentView  int
	pages        []tea.Model
}

func (m RegisterModel) Init() tea.Cmd {
	return nil
}

func (m RegisterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	focused := m.CheckFocus()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "h", "left": 
			if !focused {
				m.currentView = max(m.currentView - 1, 0);
				return m, nil
			}
		case "l", "right": 
			if !focused && m.CheckValidity() {
				m.currentView = min(m.currentView + 1, len(m.pages)-1);
				return m, nil
			}
		}
	}
	m.pages[m.currentView], cmd = m.pages[m.currentView].Update(msg)
	return m, cmd
}

func (m RegisterModel) CheckFocus() (bool) {
	view := m.pages[m.currentView]
	if v, ok := view.(RegisterOneModel); ok {
		return v.input
	} else {
		return false
	}
}

func (m RegisterModel) CheckValidity() (bool) {
	view := m.pages[m.currentView]
	if v, ok := view.(RegisterOneModel); ok {
		return v.Valid()
	} else if v, ok := view.(RegisterTwoModel); ok {
		return v.Valid()
	} else {
		return false
	}
}

func (m RegisterModel) View() string {
	var b strings.Builder
	b.WriteString(m.txtStyle.Render("Registration"))
	b.WriteString("\n")
	for i := 0; i < len(m.pages); i++ {
		if (m.currentView == i) {
			b.WriteString(m.activeDot)
		} else {
			b.WriteString(m.inactiveDot)
		}
	}
	
	b.WriteString("\n")
	b.WriteString(
		m.pageStyle.Width(35).Render(m.pages[m.currentView].View()),
	)
	b.WriteString("\n\n")
	b.WriteString(m.quitStyle.Render("h/l ←/→ page • q: quit\n"))
	return b.String()
}

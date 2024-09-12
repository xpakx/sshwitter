package main

import (
	"fmt"
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
		PageOneModel{ page: 1 },
		PageOneModel{ page: 2 },
		PageOneModel{ page: 3 },
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "h", "left": 
		        m.currentView = max(m.currentView - 1, 0);
		        return m, nil
		case "l", "right": 
		        m.currentView = min(m.currentView + 1, len(m.pages)-1);
		        return m, nil
		}
	}
	return m, cmd
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
		m.pageStyle.Width(20).Render(m.pages[m.currentView].View()),
	)
	b.WriteString("\n\n")
	b.WriteString(m.quitStyle.Render("h/l ←/→ page • q: quit\n"))
	return b.String()
}

type PageOneModel struct {
	page  int
}

func (m PageOneModel) Init() tea.Cmd {
	return nil
}

func (m PageOneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m PageOneModel) View() string {
	return fmt.Sprintf("Page %d", + m.page)
}

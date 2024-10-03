package main

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func getBoardModel(renderer *lipgloss.Renderer, db *sql.DB, user SavedUser) (BoardModel) {
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
	usernameStyle := renderer.NewStyle().Foreground(lipgloss.Color("5"))

	tabs := []tea.Model{ }
	tabs = append(tabs, getFeedView(renderer, db, user))
	tabs = append(tabs, getProfileView(renderer, db, user.username, user))


	if (user.administrator) {
		tabs = append(tabs, getModeratorTab(renderer, db))
	}


	return BoardModel{ 
		name: "sshwitter", 
		user: user,
		txtStyle: txtStyle, 
		quitStyle: quitStyle,
		userStyle: usernameStyle,
		currentTab: 0,
		tabs: tabs,
	}
}

type BoardModel struct {
	name       string
	user       SavedUser
	txtStyle   lipgloss.Style
	quitStyle  lipgloss.Style
	userStyle  lipgloss.Style
	currentTab int
	tabs       []tea.Model
}

func (m BoardModel) Init() tea.Cmd {
	return nil
}

func (m BoardModel) GetTab(index int) int {
	return min(index, len(m.tabs)-1)
}

func (m BoardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "alt+1":
			m.currentTab = m.GetTab(0)
		case "alt+2":
			m.currentTab = m.GetTab(1)
		case "alt+3":
			m.currentTab = m.GetTab(2)
		case "alt+4":
			m.currentTab = m.GetTab(3)
		case "alt+5":
			m.currentTab = m.GetTab(4)
		}
	case tea.WindowSizeMsg:
		var cmds []tea.Cmd = make([]tea.Cmd, len(m.tabs))
		for i := range m.tabs {
			m.tabs[i], cmds[i] = m.tabs[i].Update(msg)
		}
		return m, tea.Batch(cmds...)
	}
	if len(m.tabs) > 0 {
		m.tabs[m.currentTab], cmd = m.tabs[m.currentTab].Update(msg)
	}
	return m, cmd
}

func (m BoardModel) View() string {
	var tabs string
	if len(m.tabs) > 0 {
		tabs = m.tabs[m.currentTab].View()
	}
	return m.txtStyle.Render("Authorized!") + 
		"\n" + 
		"Hello, " + 
		m.userStyle.Render(m.user.username) +
		"\n\n" + 
		tabs + "\n" +
		m.quitStyle.Render("Press 'q' to quit\n")
}

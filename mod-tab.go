package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func getModeratorTab(renderer *lipgloss.Renderer) (ModeratorTabModel) {
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
	unverifiedUsers :=  GetUnverifiedUsers()
	prefixStyle := renderer.NewStyle().
			Foreground(lipgloss.Color("#1da1f2"))

	return ModeratorTabModel{ 
		txtStyle: txtStyle, 
		quitStyle: quitStyle,
		prefixStyle: prefixStyle,
		viewName: "Waiting for verification",
		users: unverifiedUsers,
		current: 0,
	}
}

func GetUnverifiedUsers() []SavedUser {
	var result []SavedUser = make([]SavedUser, 0)
	for _, user := range users {
		if !user.verified {
			result = append(result, user)
		}
	}
	return result
}

func AcceptUser(user SavedUser) {
	user.verified = true;
	users[user.username] = user
}

func DeleteUser(user SavedUser) {
	delete(users, user.username)
}

type ModeratorTabModel struct {
	txtStyle     lipgloss.Style
	quitStyle    lipgloss.Style
	prefixStyle  lipgloss.Style
	viewName     string
	users        []SavedUser
	current      int
}

func (m ModeratorTabModel) Init() tea.Cmd {
	return nil
}


func (m *ModeratorTabModel) RemoveCurrentFromList() {
	m.users = append(m.users[:m.current], m.users[m.current+1:]...)
	m.current = max(m.current - 1, 0);
}

func (m ModeratorTabModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down": 
			m.current = min(m.current + 1, len(m.users)-1);
			return m, nil
		case "k", "up": 
			m.current = max(m.current - 1, 0);
			return m, nil
		case "enter":
			if len(m.users) > 0 {
				AcceptUser(m.users[m.current])
				m.RemoveCurrentFromList()
			}
		case "delete":
			if len(m.users) > 0 {
				DeleteUser(m.users[m.current])
				m.RemoveCurrentFromList()
			}
		}
	}
	return m, nil
}

func (m ModeratorTabModel) View() string {
	doc := strings.Builder{}
	tabName := m.txtStyle.Render("Moderator tab")  
	doc.WriteString(tabName)
	doc.WriteString("\n")
	doc.WriteString(m.viewName)
	doc.WriteString("\n\n")

	if len(m.users) > 0 {
		for i, user := range m.users {
			doc.WriteString(m.getPrefix(i))
			doc.WriteString(user.username)
			doc.WriteString("\n")

		}
	} else {
		m.quitStyle.Render("No users")
	}

	return doc.String()
}

func (m ModeratorTabModel) getPrefix(i int) string {
	if m.current == i {
		return m.prefixStyle.Render("⍟ ")
	}
	return "  "
}


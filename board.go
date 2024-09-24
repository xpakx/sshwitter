package main
import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func getBoardModel(renderer *lipgloss.Renderer, user SavedUser) (BoardModel) {
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
	usernameStyle := renderer.NewStyle().Foreground(lipgloss.Color("5"))

	tabs := []tea.Model{ }
	if (user.administrator) {
		tabs = append(tabs, getModeratorTab(renderer))
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

func (m BoardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m BoardModel) View() string {
	return m.txtStyle.Render("Authorized!") + 
		"\n" + 
		"Hello, " + 
		m.userStyle.Render(m.user.username) +
		"\n\n" + 
		m.tabs[m.currentTab].View() + "\n" +
		m.quitStyle.Render("Press 'q' to quit\n")
}

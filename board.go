package main
import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func getBoardModel(renderer *lipgloss.Renderer, username string) (BoardModel) {
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
	usernameStyle := renderer.NewStyle().Foreground(lipgloss.Color("5"))

	return BoardModel{ 
		name: "sshwitter", 
		username: username,
		txtStyle: txtStyle, 
		quitStyle: quitStyle,
		userStyle: usernameStyle,
	}
}

type BoardModel struct {
	name       string
	username   string
	txtStyle   lipgloss.Style
	quitStyle  lipgloss.Style
	userStyle  lipgloss.Style
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
		m.userStyle.Render(m.username) +
		"\n\n" + 
		m.quitStyle.Render("Press 'q' to quit\n")
}

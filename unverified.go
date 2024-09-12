package main
import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func getUnverifiedModel(renderer *lipgloss.Renderer, username string) (UnverifiedModel) {
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
	usernameStyle := renderer.NewStyle().Foreground(lipgloss.Color("5"))

	return UnverifiedModel{ 
		name: "sshwitter", 
		username: username,
		txtStyle: txtStyle, 
		quitStyle: quitStyle,
		userStyle: usernameStyle,
	}
}

type UnverifiedModel struct {
	name       string
	username   string
	txtStyle   lipgloss.Style
	quitStyle  lipgloss.Style
	userStyle  lipgloss.Style
}

func (m UnverifiedModel) Init() tea.Cmd {
	return nil
}

func (m UnverifiedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m UnverifiedModel) View() string {
	return m.txtStyle.Render("Waiting for verification!") + 
		"\n" + 
		"Hello, " + 
		m.userStyle.Render(m.username) +
		". Your account is unverified." + 
		"\n\n" + 
		m.quitStyle.Render("Press 'q' to quit\n")
}

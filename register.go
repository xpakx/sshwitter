package main
import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func getRegisterModel(renderer *lipgloss.Renderer, username string) (RegisterModel) {
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
	usernameStyle := renderer.NewStyle().Foreground(lipgloss.Color("5"))

	return RegisterModel{ 
		name: "sshwitter", 
		username: username,
		txtStyle: txtStyle, 
		quitStyle: quitStyle,
		userStyle: usernameStyle,
	}
}

type RegisterModel struct {
	name       string
	username   string
	txtStyle   lipgloss.Style
	quitStyle  lipgloss.Style
	userStyle  lipgloss.Style
}

func (m RegisterModel) Init() tea.Cmd {
	return nil
}

func (m RegisterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m RegisterModel) View() string {
	return m.txtStyle.Render("Registration") + 
		"\n\n" + 
		m.quitStyle.Render("Press 'q' to quit\n")
}

package main
import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func getModeratorTab(renderer *lipgloss.Renderer) (ModeratorTabModel) {
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))

	return ModeratorTabModel{ 
		txtStyle: txtStyle, 
		quitStyle: quitStyle,
	}
}

type ModeratorTabModel struct {
	txtStyle   lipgloss.Style
	quitStyle  lipgloss.Style
}

func (m ModeratorTabModel) Init() tea.Cmd {
	return nil
}

func (m ModeratorTabModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ModeratorTabModel) View() string {
	return m.txtStyle.Render("Moderator tab")  
}

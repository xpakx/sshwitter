package main

import (
	"database/sql"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func getProfileView(renderer *lipgloss.Renderer, db *sql.DB, username string) (ProfileViewModel) {
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
	user, found :=  GetUserByUsername(db, username)
	if (!found) {
		log.Info("No such user")
		// TODO: 404 page
	}

	return ProfileViewModel{ 
		txtStyle: txtStyle, 
		quitStyle: quitStyle,
		user: user,
		db: db,
	}
}


type ProfileViewModel struct {
	txtStyle     lipgloss.Style
	quitStyle    lipgloss.Style
	user         SavedUser
	db           *sql.DB
}

func (m ProfileViewModel) Init() tea.Cmd {
	return nil
}


func (m ProfileViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ProfileViewModel) View() string {
	doc := strings.Builder{}
	profile := m.txtStyle.Render("Profile ")  
	doc.WriteString(profile)
	doc.WriteString(m.user.username)
	return doc.String()
}

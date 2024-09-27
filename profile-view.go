package main

import (
	"database/sql"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func getProfileView(renderer *lipgloss.Renderer, db *sql.DB, username string) (ProfileViewModel) {
	width := 150 // TODO
	infoWidth := 20
	infoStyle := renderer.NewStyle().
		MaxWidth(infoWidth).
		Width(infoWidth).
		PaddingRight(2)

	postStyle := renderer.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		PaddingLeft(2).
		BorderForeground(lipgloss.Color("8"))

	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))

	user, found :=  GetUserByUsername(db, username)
	if (!found) {
		log.Info("No such user")
		// TODO: 404 page
	}

	info := getProfileInfo(renderer, db, user)

	return ProfileViewModel{ 
		infoStyle: infoStyle, 
		quitStyle: quitStyle,
		postStyle: postStyle,
		infoWidth: infoWidth,
		user: user,
		db: db,
		info: info,
		width: width,
	}
}


type ProfileViewModel struct {
	infoStyle    lipgloss.Style
	quitStyle    lipgloss.Style
	postStyle    lipgloss.Style
	info         ProfileInfoModel
	user         SavedUser
	db           *sql.DB
	width        int
	infoWidth    int
}

func (m ProfileViewModel) Init() tea.Cmd {
	return nil
}


func (m ProfileViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ProfileViewModel) View() string {
	postsWidth := m.width - m.infoWidth
	info := m.infoStyle.
		Render(m.info.View())
	posts := m.postStyle.
	        Width(postsWidth).
	        MaxWidth(postsWidth).
		Render("posts")
	doc := lipgloss.JoinHorizontal(lipgloss.Left,
		info,
		posts,
	)
	return doc
}

func getProfileInfo(renderer *lipgloss.Renderer, db *sql.DB, user SavedUser) (ProfileInfoModel) {
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
	headerStyle := renderer.NewStyle().
		Foreground(lipgloss.Color("5")).
		Bold(true)
	numberStyle := renderer.NewStyle().
		Bold(true)

	return ProfileInfoModel{ 
		txtStyle: txtStyle, 
		quitStyle: quitStyle,
		headerStyle: headerStyle,
		numberStyle: numberStyle,
		user: user,
		db: db,
	}
}


type ProfileInfoModel struct {
	txtStyle     lipgloss.Style
	quitStyle    lipgloss.Style
	headerStyle  lipgloss.Style
	numberStyle  lipgloss.Style
	user         SavedUser
	db           *sql.DB
}

func (m ProfileInfoModel) Init() tea.Cmd {
	return nil
}


func (m ProfileInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ProfileInfoModel) View() string {
	doc := strings.Builder{}
	username := m.headerStyle.Render(m.user.username)  
	doc.WriteString(username)
	doc.WriteString("\n")
	description := m.txtStyle.Render("description")  
	doc.WriteString(description)
	doc.WriteString("\n\n")
	doc.WriteString("📍 " )
	location := m.quitStyle.Render("City")  
	doc.WriteString(location)
	doc.WriteString("\n")
	doc.WriteString("🗓  " )
	joinDate := m.quitStyle.Render("Date")  
	doc.WriteString(joinDate)
	doc.WriteString("\n")
	doc.WriteString(m.numberStyle.Render("0"))
	doc.WriteString(m.quitStyle.Render(" Following"))
	doc.WriteString("\n")
	doc.WriteString(m.numberStyle.Render("0"))
	doc.WriteString(m.quitStyle.Render(" Followers"))
	return doc.String()
}

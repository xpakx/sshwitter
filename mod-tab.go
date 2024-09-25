package main

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func getModeratorTab(renderer *lipgloss.Renderer, db *sql.DB) (ModeratorTabModel) {
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
	unverifiedUsers, err :=  GetUnverifiedUsers(db)
	if (err != nil) {
		log.Debug("Error while fetching users")
	}


	columns := []table.Column{
		{Title: "Id", Width: 3},
		{Title: "Username", Width: 10},
		{Title: "Email", Width: 10},
		{Title: "Birth date", Width: 10},
	}

	var rows []table.Row = make([]table.Row, 0, len(unverifiedUsers))
	for _, user  := range unverifiedUsers {
		rows = append(rows, table.Row{strconv.Itoa(int(user.id)), user.username, user.email, ""})
		
	}
	table := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)


	prefixStyle := renderer.NewStyle().
			Foreground(lipgloss.Color("#1da1f2"))

	return ModeratorTabModel{ 
		txtStyle: txtStyle, 
		quitStyle: quitStyle,
		prefixStyle: prefixStyle,
		viewName: "Waiting for verification",
		users: unverifiedUsers,
		current: 0,
		db: db,
		table: table,
	}
}


type ModeratorTabModel struct {
	txtStyle     lipgloss.Style
	quitStyle    lipgloss.Style
	prefixStyle  lipgloss.Style
	viewName     string
	users        []SavedUser
	current      int
	db           *sql.DB
	table        table.Model
}

func (m ModeratorTabModel) Init() tea.Cmd {
	return nil
}

func (m ModeratorTabModel) GetCurrentIndex() int {
	username := m.table.SelectedRow()[1]
	for i, user := range m.users {
		if username == user.username {
			return i
		}
	}
	return -1
}


func (m *ModeratorTabModel) RemoveCurrentFromList() {
	current := m.GetCurrentIndex()
	if current < 0 {
		return
	}
	m.users = append(m.users[:current], m.users[current+1:]...)
	var rows []table.Row = make([]table.Row, 0, len(m.users))
	for _, user  := range m.users {
		rows = append(rows, table.Row{strconv.Itoa(int(user.id)), user.username, user.email, ""})
	}
	m.table.SetRows(rows)
}

func (m ModeratorTabModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if len(m.users) > 0 {
				// TODO: move to command
				user := m.GetCurrentChoice()
				AcceptUser(m.db, user)
				m.RemoveCurrentFromList()
			}
		case "delete":
			if len(m.users) > 0 {
				// TODO: move to command
				user := m.GetCurrentChoice()
				DeleteUser(m.db, user)
				m.RemoveCurrentFromList()
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m ModeratorTabModel) GetCurrentChoice() SavedUser {
	username := m.table.SelectedRow()[1]
	for _, user := range m.users {
		if username == user.username {
			return user
		}
	}
	return SavedUser{}
}

func (m ModeratorTabModel) View() string {
	doc := strings.Builder{}
	tabName := m.txtStyle.Render("Moderator tab")  
	doc.WriteString(tabName)
	doc.WriteString("\n")
	doc.WriteString(m.viewName)
	doc.WriteString("\n\n")

	if len(m.users) > 0 {
		doc.WriteString(m.table.View())
	} else {
		doc.WriteString(m.quitStyle.Render("No users"))
	}


	return doc.String()
}

func (m ModeratorTabModel) getPrefix(i int) string {
	if m.current == i {
		return m.prefixStyle.Render("⍟ ")
	}
	return "  "
}


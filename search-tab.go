package main

import (
	"database/sql"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func getSearchView(renderer *lipgloss.Renderer, db *sql.DB, user SavedUser) (Tab) {
	infoWidth := 20
	infoStyle := renderer.NewStyle().
		MaxWidth(infoWidth).
		Width(infoWidth).
		PaddingRight(2).
		PaddingLeft(1)

	postStyle := renderer.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		PaddingLeft(2).
		BorderForeground(lipgloss.Color("8"))

	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))

	headerStyle := renderer.NewStyle().
		Foreground(lipgloss.Color("5"))
	numberStyle := quitStyle.
		Bold(true)

	nameInput := CreateCustomInput(renderer, "Search", "name", nameValidator, true)


	return Tab{
		Model: SearchViewModel{ 
			infoStyle: infoStyle, 
			quitStyle: quitStyle,
			postStyle: postStyle,
			infoWidth: infoWidth,
			headerStyle: headerStyle,
			numberStyle: numberStyle,
			db: db,
			renderer: renderer,
			user: user,
			nameInput: nameInput,
			input: false,
		},
		Name: "Search",
	}
}


type SearchViewModel struct {
	infoStyle    lipgloss.Style
	quitStyle    lipgloss.Style
	postStyle    lipgloss.Style
	headerStyle  lipgloss.Style
	numberStyle  lipgloss.Style
	user         SavedUser
	db           *sql.DB
	renderer     *lipgloss.Renderer
	width        int
	infoWidth    int
	nameInput    CustomInput
	input        bool
	users        []SavedUser
}

func (m SearchViewModel) Init() tea.Cmd {
	return nil
}

func (m SearchViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: 
		m.width = msg.Width
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "p":
			if (!m.input) {
				return m, nil
			}
		case "esc":
			m.nameInput.Blur()
			m.input = false
			return m, nil
		case "enter": 
		        if(!m.input) {
				m.input = true
				return m, m.nameInput.Focus()
			} else {
				m.nameInput.Blur()
				searchQuery := m.nameInput.Input.Value()
				m.input = false;
				users, err := SearchUsers(m.db, searchQuery)
				if (err == nil) {
					m.users =  users
				} else {
					log.Info(err)
				}
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.nameInput, cmd = m.nameInput.Update(msg)
	return m, cmd
}


func (m SearchViewModel) View() string {
	doc := strings.Builder{}
	name := m.nameInput.View(false)
	doc.WriteString(name)
	if (len(m.users) != 0) {
		for _, user := range m.users {
			doc.WriteString("\n")
			doc.WriteString(user.username)
		}
	}
	return doc.String()
}

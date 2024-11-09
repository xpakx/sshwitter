package main

import (
	"database/sql"
	"fmt"
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

	nameInput := CreateCustomInput(renderer, "Search", "name", searchQueryValidator, true)


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
			current: 0,
			inList: false,
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
	current      int
	inList       bool
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
		case "esc":
			m.nameInput.Blur()
			m.input = false
			return m, nil
		case "enter": 
		        if(!m.input) {
				if (!m.inList) {
					m.input = true
					return m, m.nameInput.Focus()
				} else {
					if (m.current == 0) {
						return m, nil
					}
					curr := m.users[m.current-1]
					return m, openProfile(curr.username)
				}
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
		case "j", "down": 
		        if(!m.input) {
				m.current = min(m.current + 1, len(m.users));
				if (m.current > 0) {
					m.inList = true
				}
				return m, nil
			}
		case "k", "up": 
		        if(!m.input) {
				m.current = max(m.current - 1, 0);
				if (m.current == 0) {
					m.inList = false
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
		doc.WriteString("\n")
		doc.WriteString(m.quitStyle.Render("Results"))
		for i, user := range m.users {
			doc.WriteString("\n")
			if (m.inList && m.current-1 == i) {
				doc.WriteString("*")
			} else {
				doc.WriteString(" ")
			}
			doc.WriteString(user.username)
		}
	} else {
		doc.WriteString("\n")
		doc.WriteString(m.quitStyle.Render("No results"))
	}
	return doc.String()
}

func searchQueryValidator(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("cannot be empty")
	}
	return nil
}

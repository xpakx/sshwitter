package main

import (
	"database/sql"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EditProfileModel struct {
	descriptionInput CustomInput
	locationInput  CustomInput
	elems          int
	current        int
	err            error
	input          bool
	headerStyle    lipgloss.Style
	subheaderStyle lipgloss.Style
	buttonStyle    lipgloss.Style
	db             *sql.DB
	user           SavedUser
}

func getEditProfileModel(renderer *lipgloss.Renderer, db *sql.DB, user SavedUser) Tab {
	descriptionInput := CreateCustomInput(renderer, "Description", "Describe yourself", descriptionValidator, true)
	if user.description.Valid {
		descriptionInput.Input.SetValue(user.description.String)
	}
	locationInput := CreateCustomInput(renderer, "Location", "Where are you?", locationValidator, false)
	if user.location.Valid {
		locationInput.Input.SetValue(user.location.String)
	}

	headerStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	subheaderStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))

	buttonStyle := renderer.NewStyle().
		MarginTop(2).
		Width(30).
		Align(lipgloss.Right)

	return Tab{
		Model: EditProfileModel{
			descriptionInput: descriptionInput,
			locationInput:    locationInput,
			elems:            3,
			err:              nil,
			input:            true,
			headerStyle:      headerStyle,
			subheaderStyle:   subheaderStyle,
			buttonStyle:      buttonStyle,
			db:               db,
			user:             user,
		},
		Name: "Edit profile",
	}
}

func (m EditProfileModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m EditProfileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, 2)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.descriptionInput.Blur()
			m.locationInput.Blur()
			m.input = false
			return m, nil
		case "j", "down":
			if !m.input {
				m.current = min(m.current+1, m.elems-1)
				return m, nil
			}
		case "k", "up":
			if (!m.input) {
				m.current = max(m.current-1, 0)
				return m, nil
			}
		case "enter":
			if !m.input {
				if m.current == 0 {
					m.input = true
					return m, m.descriptionInput.Focus()
				} else if m.current == 1 {
					m.input = true
					return m, m.locationInput.Focus()
				} else if m.current == 2 {
					if !m.Valid() {
						return m, nil
					}
					desc := m.descriptionInput.Input.Value()
					loc := m.locationInput.Input.Value()
					
					err := UpdateUserData(m.db, m.user, desc, loc)
					if err != nil {
						return m, nil
					}
					return m, closeEdit(desc, loc)
				}
			} else {
				m.descriptionInput.Blur()
				m.locationInput.Blur()
				m.input = false
			}
		}
	case error:
		m.err = msg
		return m, nil
	}

	m.descriptionInput, cmds[0] = m.descriptionInput.Update(msg)
	m.locationInput, cmds[1] = m.locationInput.Update(msg)
	return m, tea.Batch(cmds...)
}

func (m EditProfileModel) View() string {
	description := m.descriptionInput.View(m.current == 0)
	location := m.locationInput.View(m.current == 1)

	var button string

	if m.current == 2 {
		button = m.buttonStyle.
			Render(getButtonPrefix(m.current == 2) + "[ Save ]")
	} else {
		button = m.buttonStyle.
			Foreground(lipgloss.Color("8")).
			Render("[ Save ]")
	}

	return m.headerStyle.Render("Edit your profile") + "\n" +
		"\n\n" +
		description +
		"\n" +
		location +
		"\n" +
		button
}

func (m EditProfileModel) Valid() bool {
	return m.descriptionInput.Valid() &&
		m.locationInput.Valid() &&
		len(m.descriptionInput.Input.Value()) > 0 &&
		len(m.locationInput.Input.Value()) > 0
}

func descriptionValidator(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("cannot be empty")
	}
	if len(s) > 100 {
		return fmt.Errorf("description is too long")
	}
	return nil
}

func locationValidator(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("cannot be empty")
	}
	if len(s) > 50 {
		return fmt.Errorf("location is too long")
	}
	return nil
}

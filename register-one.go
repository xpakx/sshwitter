package main

import (
	"fmt"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	//"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textinput"
)

type RegisterOneModel struct {
	page       int
	steps       int
	nameInput  textinput.Model
	emailInput textinput.Model
	birthInput textinput.Model
	elems      int
	current    int
	err        error
	input      bool
}

func getPageOneModel(steps int) RegisterOneModel {
	nameInput := textinput.New()
	nameInput.Placeholder = "Name"
	nameInput.Focus()
	nameInput.CharLimit = 40
	nameInput.Width = 20
	nameInput.Validate = nameValidator

	emailInput := textinput.New()
	emailInput.Placeholder = "Mail"
	emailInput.CharLimit = 40
	emailInput.Width = 20
	emailInput.Validate = emailValidator

	birthInput := textinput.New()
	birthInput.Placeholder = "yyyy-mm-dd"
	birthInput.CharLimit = 40
	birthInput.Width = 20
	birthInput.Validate = dateValidator

	return RegisterOneModel {
		nameInput: nameInput,
		emailInput: emailInput,
		birthInput: birthInput,
		page: 1,
		steps: steps,
		elems: 3,
		err:       nil,
		input: true,
	}
}

func (m RegisterOneModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m RegisterOneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, 3)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.nameInput.Blur()
			m.emailInput.Blur()
			m.birthInput.Blur()
			m.input = false
			return m, nil;
		case "j", "down": 
		        if(!m.input) {
				m.current = min(m.current + 1, m.elems-1);
				return m, nil
			}
		case "k", "up": 
		        if(!m.input) {
				m.current = max(m.current - 1, 0);
				return m, nil
			}
		case "enter": 
		        if(!m.input) {
				m.input = true
				if (m.current == 0) {
					return m, m.nameInput.Focus()
				} else if (m.current == 1) {
					m.emailInput.Focus()
					return m, m.emailInput.Focus()
				} else if (m.current == 2) {
					m.birthInput.Focus()
					return m, m.birthInput.Focus()
				}
			} else {
				m.nameInput.Blur()
				m.emailInput.Blur()
				m.birthInput.Blur()
				m.input = false;
			}
		}
	case error:
		m.err = msg
		return m, nil
	}

	m.nameInput, cmds[0] = m.nameInput.Update(msg)
	m.emailInput, cmds[1] = m.emailInput.Update(msg)
	m.birthInput, cmds[2] = m.birthInput.Update(msg)
	return m, tea.Batch(cmds...)
}

func getInputPrefix(input textinput.Model, current bool) string {
	if current {
		if (input.Err != nil) {
			return "⨯ "
		} else {
			return "o "
		}
	}
	return "  "
}

func (m RegisterOneModel) View() string {
	preName := getInputPrefix(m.nameInput, m.current == 0)
	preEmail := getInputPrefix(m.emailInput, m.current == 1)
	preBirth := getInputPrefix(m.birthInput, m.current == 2)
	return "Create your account\n" +
		fmt.Sprintf("Step %d of %d", m.page, m.steps) +
		"\n\n" +
		preName + m.nameInput.View() +
		"\n" +
		preEmail + m.emailInput.View() +
		"\n" +
		preBirth + m.birthInput.View()

}

func nameValidator(s string) error {
	if (len(s) > 20 || len(s) < 5) {
		return fmt.Errorf("Name should be between 5 and 20 characters")
	}
	return nil
}

func emailValidator(s string) error {
	if (!strings.Contains(s, "@")) {
		return fmt.Errorf("Should be correct email")
	}
	return nil
}

func dateValidator(s string) error {
	match, _ := regexp.MatchString("^[0-9]{4}-[0-9]{2}-[0-9]{2}$", s)
	if (!match) {
		return fmt.Errorf("Date should have format yyyy-mm-dd")
	}
	return nil
}

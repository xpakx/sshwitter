package main

import (
	"fmt"
	//"strings"

	tea "github.com/charmbracelet/bubbletea"
	//"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textinput"
)

type RegisterOneModel struct {
	page       int
	steps       int
	nameInput  textinput.Model
	emailInput textinput.Model
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

	emailInput := textinput.New()
	emailInput.Placeholder = "Mail"
	emailInput.CharLimit = 40
	emailInput.Width = 20

	return RegisterOneModel {
		nameInput: nameInput,
		emailInput: emailInput,
		page: 1,
		steps: steps,
		elems: 2,
		err:       nil,
		input: true,
	}
}

func (m RegisterOneModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m RegisterOneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, 2)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.nameInput.Blur()
			m.emailInput.Blur()
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
				}
			} else {
				m.nameInput.Blur()
				m.emailInput.Blur()
				m.input = false;
			}
		}
	case error:
		m.err = msg
		return m, nil
	}

	m.nameInput, cmds[0] = m.nameInput.Update(msg)
	m.emailInput, cmds[1] = m.emailInput.Update(msg)
	return m, tea.Batch(cmds...)
}

func (m RegisterOneModel) View() string {
	preName := "  " 
	if m.current == 0 {
		preName = "o "
	}
	preEmail := "  " 
	if m.current == 1 {
		preEmail = "o "
	}
	return "Create your account\n" +
		fmt.Sprintf("Step %d of %d", m.page, m.steps) +
		"\n\n" +
		preName + m.nameInput.View() +
		"\n" +
		preEmail + m.emailInput.View()

}


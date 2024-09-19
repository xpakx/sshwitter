package main

import (
	"fmt"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textinput"
)

type RegisterOneModel struct {
	page       int
	steps      int
	nameInput  CustomInput
	emailInput CustomInput
	birthInput CustomInput
	elems      int
	current    int
	err        error
	input      bool
}

func getPageOneModel(steps int) RegisterOneModel {
	nameInput := createCustomInput("User name", "name", nameValidator, true)
	emailInput := createCustomInput("E-mail", "mail", emailValidator, false)
	birthInput := createCustomInput("Birth date", "yyyy-mm-dd", dateValidator, false)

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
	name := m.nameInput.View(m.current == 0)
	email := m.emailInput.View(m.current == 1)
	birth := m.birthInput.View(m.current == 2)
	return "Create your account\n" +
		fmt.Sprintf("Step %d of %d", m.page, m.steps) +
		"\n\n" +
		name +
		"\n" +
		email +
		"\n" +
		birth

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




type CustomInput struct {
	Name       string
	Input      textinput.Model
}

func createCustomInput(name string, placeholder string, validator textinput.ValidateFunc, autofocus bool) CustomInput {
	input := textinput.New()
	input.Placeholder = placeholder
	input.Prompt  = ""

	input.CharLimit = 40
	input.Width = 25
	input.Validate = validator

	if autofocus {
		input.Focus()
	}

	return CustomInput{
		Name: name,
		Input: input,
	}
}

func (i CustomInput) Valid() bool {
	return i.Input.Err == nil
}

func (i CustomInput) Invalid() bool {
	return i.Input.Err != nil
}

func (i *CustomInput) Blur() {
	i.Input.Blur()
}

func (i CustomInput) getPrefix(current bool) string {
	if current {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#1da1f2")).
			Render("⍟ ")
	}
	return "  "
}


func (i CustomInput) getBorderColor(current bool) lipgloss.Color {
	if current && i.Input.Focused() {
		return lipgloss.Color("#1da1f2")
	}
	if i.Invalid() {
		return lipgloss.Color("#cc0033")
	}
	return lipgloss.Color("8")
}

func (i CustomInput) getValidationPrefix() string {
		if (i.Invalid()) {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("#cc0033")).
				Render("⨯ ")
		} else if (i.Valid() && len(i.Input.Value()) > 0) {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("#34b233")).
				Render("✔ ")
		} else {
			return "  "
		}
}

func (i CustomInput) View(current bool) string {
	nameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("5")).
		Bold(true)

	inputStyle := lipgloss.NewStyle().
	        Border(lipgloss.RoundedBorder())


	error := "";
	if(i.Invalid()) {
		error = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#cc0033")).
			PaddingLeft(2).
			MaxHeight(1).
			MaxWidth(25).
			Render(i.Input.Err.Error())
	}

	output := lipgloss.JoinVertical(
		lipgloss.Top,
		nameStyle.Render(i.Name),
		lipgloss.JoinHorizontal(lipgloss.Left,
			"\n" + i.getValidationPrefix() + "\n",
			inputStyle.BorderForeground(i.getBorderColor(current)).Render(i.Input.View()),
		),
		error,
	)
	return i.getPrefix(current) + output
}

func (i CustomInput) Update(msg tea.Msg) (CustomInput, tea.Cmd) {
	var cmd tea.Cmd
	i.Input, cmd = i.Input.Update(msg)
	return i, cmd
}

func (i *CustomInput) Focus() tea.Cmd {
	return i.Input.Focus()
}

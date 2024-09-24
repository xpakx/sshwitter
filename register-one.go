package main

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RegisterOneModel struct {
	page             int
	steps            int
	nameInput        CustomInput
	emailInput       CustomInput
	birthInput       CustomInput
	elems            int
	current          int
	err              error
	input            bool
	headerStyle      lipgloss.Style
	subheaderStyle   lipgloss.Style
	buttonStyle      lipgloss.Style
	
}

func getPageOneModel(renderer *lipgloss.Renderer, steps int, username string) RegisterOneModel {
	nameInput := createCustomInput(renderer, "User name", "name", nameValidator, true)
	nameInput.Input.SetValue(username)
	emailInput := createCustomInput(renderer, "E-mail", "mail", emailValidator, false)
	birthInput := createCustomInput(renderer, "Birth date", "yyyy-mm-dd", dateValidator, false)


	headerStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	subheaderStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))

	buttonStyle := lipgloss.NewStyle().
		MarginTop(2).
		Width(30).
		Align(lipgloss.Right)

	return RegisterOneModel {
		nameInput: nameInput,
		emailInput: emailInput,
		birthInput: birthInput,
		page: 1,
		steps: steps,
		elems: 4,
		err:       nil,
		input: true,
		headerStyle: headerStyle,
		subheaderStyle: subheaderStyle,
		buttonStyle: buttonStyle,
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
				if (m.current == 0) {
					m.input = true
					return m, m.nameInput.Focus()
				} else if (m.current == 1) {
					m.input = true
					m.emailInput.Focus()
					return m, m.emailInput.Focus()
				} else if (m.current == 2) {
					m.input = true
					m.birthInput.Focus()
					return m, m.birthInput.Focus()
				} else if (m.current == 3) {
					return m, nextPage
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

	var button string;
	
	if (m.current == 3)  {
		button = m.buttonStyle.
			Render(getButtonPrefix(m.current == 3) + "[ Next ]")
	} else {
		button = m.buttonStyle.
			Foreground(lipgloss.Color("8")).
			Render("[ Next ]")
	}

	return m.headerStyle.Render("Create your account") + "\n" +
		m.subheaderStyle.Render(fmt.Sprintf("Step %d of %d", m.page, m.steps)) +
		"\n\n" +
		name +
		"\n" +
		email +
		"\n" +
		birth +
		button
}

func (m RegisterOneModel) Valid() bool {
	return m.nameInput.Valid() &&
		m.emailInput.Valid() &&
		m.birthInput.Valid() &&
		len(m.nameInput.Input.Value()) > 0 &&
		len(m.emailInput.Input.Value()) > 0 &&
		len(m.birthInput.Input.Value()) > 0
}

func nameValidator(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("what's your name?")
	}
	if (len(s) < 5) {
		return fmt.Errorf("at least 5 characters")
	}
	if (len(s) > 20) {
		return fmt.Errorf("at most 20 characters")
	}
	return nil
}

func emailValidator(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("please enter email")
	}
	_, err := mail.ParseAddress(s)
	if (err != nil) {
		return fmt.Errorf("should be correct email")
	}
	return nil
}

func dateValidator(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("please enter birthday")
	}
	_, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("format: yyyy-mm-dd")
	}
	return nil
}




type CustomInput struct {
	Name             string
	Input            textinput.Model
	prefixStyle      lipgloss.Style
	errorMarkStyle   lipgloss.Style
	checkmarkStyle   lipgloss.Style
	inputStyle       lipgloss.Style
	nameStyle        lipgloss.Style
	errorStyle       lipgloss.Style
}

func createCustomInput(renderer *lipgloss.Renderer, name string, placeholder string, validator textinput.ValidateFunc, autofocus bool) CustomInput {
	input := textinput.New()
	input.Placeholder = placeholder
	input.Prompt  = ""

	input.CharLimit = 40
	input.Width = 25
	input.Validate = validator

	prefixStyle := renderer.NewStyle().
			Foreground(lipgloss.Color("#1da1f2"))
	errorMarkStyle := renderer.NewStyle().
			Foreground(lipgloss.Color("#cc0000"))
	checkmarkStyle := renderer.NewStyle().
			Foreground(lipgloss.Color("#34b233"))
	nameStyle := renderer.NewStyle().
			Foreground(lipgloss.Color("5")).
			Bold(true)
	errorStyle :=  renderer.NewStyle().
			Foreground(lipgloss.Color("#cc0000")).
			PaddingLeft(2).
			MaxHeight(1).
			MaxWidth(25)

	inputStyle := renderer.NewStyle().
			Border(lipgloss.RoundedBorder())

	if autofocus {
		input.Focus()
	}

	return CustomInput{
		Name: name,
		Input: input,
		prefixStyle: prefixStyle,
		errorMarkStyle: errorMarkStyle,
		checkmarkStyle: checkmarkStyle,
		nameStyle: nameStyle,
		inputStyle: inputStyle,
		errorStyle: errorStyle,
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
		return i.prefixStyle.Render("⍟ ")
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
			return i.errorMarkStyle.Render("⨯ ")
		} else if (i.Valid() && len(i.Input.Value()) > 0) {
			return i.checkmarkStyle.Render("✔ ")
		} else {
			return "  "
		}
}

func (i CustomInput) View(current bool) string {
	error := "";
	if(i.Invalid()) {
		error = i.errorStyle.Render(i.Input.Err.Error())
	}

	output := lipgloss.JoinVertical(
		lipgloss.Top,
		i.nameStyle.Render(i.Name),
		lipgloss.JoinHorizontal(lipgloss.Left,
			"\n" + i.getValidationPrefix() + "\n",
			i.inputStyle.BorderForeground(i.getBorderColor(current)).Render(i.Input.View()),
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

package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


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

func CreateCustomInput(renderer *lipgloss.Renderer, name string, placeholder string, validator textinput.ValidateFunc, autofocus bool) CustomInput {
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

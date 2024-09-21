package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RegisterTwoModel struct {
	page             int
	steps            int
	elems            int
	current          int
	err              error
	input            bool
	accepted         bool
	headerStyle      lipgloss.Style
	subheaderStyle   lipgloss.Style
}

func getPageTwoModel(steps int) RegisterTwoModel {
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	subheaderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	return RegisterTwoModel {
		page: 2,
		steps: steps,
		elems: 3, // 1-terms, 2-checkbox, 3-next button
		err:       nil,
		input: false,
		accepted: false,
		headerStyle: headerStyle,
		subheaderStyle: subheaderStyle,
		current: 0,
	}
}

func (m RegisterTwoModel) Init() tea.Cmd {
	return nil
}

type NextPageMsg struct {}

func nextPage() tea.Msg {
	return NextPageMsg{}
}

func (m RegisterTwoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd 

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, nil;
		case "j", "down": 
			m.current = min(m.current + 1, m.elems-1);
			return m, nil
		case "k", "up": 
			m.current = max(m.current - 1, 0);
			return m, nil
		case "enter": 
		        if(m.current == 1) {
				m.accepted = !m.accepted
				return m, nil
			} else if(m.current == 2) {
				return m, nextPage
			}
		}
	case error:
		m.err = msg
		return m, nil
	}

	return m, cmd
}

func (m RegisterTwoModel) RenderAccepted() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#34b233"))
	if m.accepted {
		return "[" + style.Render("✔") + "]"
	}
	return "[ ]"
}


func getButtonPrefix(current bool) string {
	if current {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#1da1f2")).
			Render("⍟ ")
	}
	return "  ";
}

func (m RegisterTwoModel) View() string {
	accepted := lipgloss.JoinHorizontal(lipgloss.Top,
		getButtonPrefix(m.current == 1) + m.RenderAccepted(),
		 " Agree to terms",
		)

	buttonStyle := lipgloss.NewStyle().
		MarginTop(2).
		Width(30).
		Align(lipgloss.Right)
	var button string;
	
	if (m.current == 2)  {
		button = buttonStyle.
			Render(getButtonPrefix(m.current == 2) + "[ Next ]")
	} else {
		button = buttonStyle.
			Foreground(lipgloss.Color("8")).
			Render("[ Next ]")
	}

	return m.headerStyle.Render("Accept terms") + "\n" +
		m.subheaderStyle.Render(fmt.Sprintf("Step %d of %d", m.page, m.steps)) +
		"\n\n" +
		"Blah blah blah\n\n" +
		accepted +
		button

}

func (m RegisterTwoModel) Valid() bool {
	return m.accepted
}


package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type RegisterTwoModel struct {
	page       int
	steps      int
	elems      int
	current    int
	err        error
	input      bool
	accepted   bool
}

func getPageTwoModel(steps int) RegisterTwoModel {
	return RegisterTwoModel {
		page: 2,
		steps: steps,
		elems: 3, // 1-terms, 2-checkbox, 3-next button
		err:       nil,
		input: false,
		accepted: false,
	}
}

func (m RegisterTwoModel) Init() tea.Cmd {
	return nil
}

func (m RegisterTwoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd 

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
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
			}
		}
	case error:
		m.err = msg
		return m, nil
	}

	return m, cmd
}

func (m RegisterTwoModel) View() string {
	var accepted string
	if (m.accepted) {
		accepted = "[x]"
	} else {
		accepted = "[ ]"
	}
	return "Accept terms\n" +
		fmt.Sprintf("Step %d of %d", m.page, m.steps)  +
		"\n\n" +
		"Blah blah blah\n\n" +
		accepted + " Agree to terms\n"  +
		"[Next]"

}

func (m RegisterTwoModel) Valid() bool {
	return m.accepted
}


package main

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Tab struct {
    tea.Model
    Name      string
}

func getBoardModel(renderer *lipgloss.Renderer, db *sql.DB, user SavedUser) (BoardModel) {
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
	usernameStyle := renderer.NewStyle().Foreground(lipgloss.Color("5"))

	tabs := []Tab{ }
	tabs = append(tabs, getFeedView(renderer, db, user, FindAllPosts, "Feed"))
	tabs = append(tabs, getFeedView(renderer, db, user, FindFollowedPosts, "Follows"))
	tabs = append(tabs, getFeedView(renderer, db, user, FindLikedPosts, "Likes"))
	tabs = append(tabs, getProfileView(renderer, db, user.username, user))


	if (user.administrator) {
		tabs = append(tabs, getModeratorTab(renderer, db))
	}

	activeTabBorder := lipgloss.Border{
		Top:         "─", Bottom:      " ", Left:        "│",
		Right:       "│", TopLeft:     "╭", TopRight:    "╮",
		BottomLeft:  "┘", BottomRight: "└",
	}

	tabBorder := lipgloss.Border{
		Top:         "─", Bottom:      "─", Left:        "│",
		Right:       "│", TopLeft:     "╭", TopRight:    "╮",
		BottomLeft:  "┴", BottomRight: "┴",
	}

	tabStyle := lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(lipgloss.Color("10")).
		Foreground(lipgloss.Color("8")).
		Padding(0, 1)

	activeTabStyle := lipgloss.NewStyle().
		Border(activeTabBorder, true).
		BorderForeground(lipgloss.Color("10")).
		Padding(0, 1)

	return BoardModel{ 
		name: "sshwitter", 
		user: user,
		txtStyle: txtStyle, 
		quitStyle: quitStyle,
		userStyle: usernameStyle,
		currentTab: 0,
		tabs: tabs,
		tabStyle: tabStyle,
		aTabStyle: activeTabStyle,
		renderer: renderer,
		db: db,
	}
}

type BoardModel struct {
	name       string
	user       SavedUser
	txtStyle   lipgloss.Style
	quitStyle  lipgloss.Style
	userStyle  lipgloss.Style
	tabStyle   lipgloss.Style
	aTabStyle  lipgloss.Style
	currentTab int
	tabs       []Tab
	renderer   *lipgloss.Renderer
	db         *sql.DB
	lastResize tea.WindowSizeMsg
}

func (m BoardModel) Init() tea.Cmd {
	return nil
}

func (m BoardModel) GetTab(index int) int {
	return min(index, len(m.tabs)-1)
}

func (m BoardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "alt+1", "alt+2", "alt+3", "alt+4", "alt+5":
			tabNumber := int(msg.String()[4] - '0')
			m.currentTab = m.GetTab(tabNumber - 1) 
			return m, nil
		case "alt+x":
			return m, closeTab(m.currentTab)
		case "alt+a":
			return m, openFeed(allFeed)
		case "alt+f":
			return m, openFeed(followedFeed)
		case "alt+l":
			return m, openFeed(likedFeed)
		case "alt+H":
			return m, openHome
		}
	case tea.WindowSizeMsg:
		m.lastResize = msg
		var cmds []tea.Cmd = make([]tea.Cmd, len(m.tabs))
		for i := range m.tabs {
			m.tabs[i].Model, cmds[i] = m.tabs[i].Model.Update(msg)
		}
		return m, tea.Batch(cmds...)
	case CloseTabMsg:
		if len(m.tabs) == 0 {
			return m, nil
		}
		i := msg.page
		m.tabs = append(m.tabs[:i], m.tabs[i+1:]...)
		if i == len(m.tabs) && i != 0 {
			m.currentTab = m.currentTab - 1
		}
		return m, nil
	case OpenFeedMsg:
		for _, tab  := range m.tabs {
			if tab.Name == msg.name {
				return m, nil
			}
		}
		m.tabs = append(
			m.tabs, 
			getFeedView(m.renderer, m.db, m.user, msg.find, msg.name),
		)
		m.tabs[len(m.tabs)-1].Model, cmd = m.tabs[len(m.tabs)-1].Model.Update(m.lastResize)
		return m, cmd
	case OpenHomeMsg:
		for _, tab  := range m.tabs {
			if tab.Name == m.user.username {
				return m, nil
			}
		}
		m.tabs = append(m.tabs, getProfileView(m.renderer, m.db, m.user.username, m.user))
		m.tabs[len(m.tabs)-1].Model, cmd = m.tabs[len(m.tabs)-1].Model.Update(m.lastResize)
		return m, nil
	}
	if len(m.tabs) > 0 {
		m.tabs[m.currentTab].Model, cmd = m.tabs[m.currentTab].Model.Update(msg)
	}
	return m, cmd
}

func (m BoardModel) View() string {
	var currentTab string
	if len(m.tabs) > 0 {
		currentTab = m.tabs[m.currentTab].View()
	}
	info := m.quitStyle.Render("Press 'q' to quit\n")

	var tabs []string = make([]string, len(m.tabs))
	for i, tab := range m.tabs {
		if i == m.currentTab {
			tabs[i] = m.aTabStyle.Render(tab.Name)
		} else {
			tabs[i] = m.tabStyle.Render(tab.Name)
		}
		
	}
	row := lipgloss.JoinHorizontal(lipgloss.Left, tabs...)
	row = lipgloss.JoinVertical(lipgloss.Top, row, currentTab, info)
	return row
}

type CloseTabMsg struct {
	page int
}

type OpenFeedMsg struct {
	name string
	find FindPostsFunc
}

func closeTab(page int) tea.Cmd {
	return func() tea.Msg {
		return CloseTabMsg{page: page}
	}
}

type FeedType int

const (
        allFeed = iota
        followedFeed = iota
        likedFeed = iota
)

func openFeed(feed FeedType) tea.Cmd {
	return func() tea.Msg {
		switch (feed) {
			case allFeed: return OpenFeedMsg{name: "Feed", find: FindAllPosts};
			case followedFeed: return OpenFeedMsg{name: "Follows", find: FindFollowedPosts};
			case likedFeed: return OpenFeedMsg{name: "Likes", find: FindLikedPosts};
			default: return OpenFeedMsg{name: "Feed", find: FindAllPosts};
		}
	}
}

type OpenHomeMsg struct {}

func openHome() tea.Msg {
	return OpenHomeMsg{}
}

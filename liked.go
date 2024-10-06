package main

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/bubbles/viewport"
)

func getLikedFeedView(renderer *lipgloss.Renderer, db *sql.DB, user SavedUser) (LikedFeedModel) {
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

	posts, err := FindLikedPosts(db, user)
	if err != nil {
		log.Error(err)
	}
	timeline := getTimeline(renderer, db, posts, user)

	newViewport := viewport.New(20, 15)

	return LikedFeedModel{ 
		infoStyle: infoStyle, 
		quitStyle: quitStyle,
		postStyle: postStyle,
		headerStyle: headerStyle,
		numberStyle: numberStyle,
		db: db,
		renderer: renderer,
		posts: timeline,
		user: user,
		viewport: newViewport,
	}
}


type LikedFeedModel struct {
	infoStyle    lipgloss.Style
	quitStyle    lipgloss.Style
	postStyle    lipgloss.Style
	headerStyle  lipgloss.Style
	numberStyle  lipgloss.Style
	posts        TimelineModel
	user         SavedUser
	db           *sql.DB
	renderer     *lipgloss.Renderer
	width        int
	viewport     viewport.Model
}

func (m LikedFeedModel) Init() tea.Cmd {
	return nil
}

func (m LikedFeedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: 
		m.width = msg.Width
		m.posts.width = max(m.width, 20) - 2
		m.viewport.Width = m.posts.width
		m.viewport.Height = msg.Height - 5
		m.viewport.SetContent(m.posts.View())
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			posts, err := FindAllPosts(m.db, m.user)
			if err != nil {
				log.Error(err)
			}
			m.posts = getTimeline(m.renderer, m.db, posts, m.user)
			m.viewport.SetContent(m.posts.View())
			return m, nil
		}
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m LikedFeedModel) View() string {
	postsWidth := max(m.width, 20)
	posts := make([]string, 0)
	posts = append(posts, m.viewport.View())
	renderedPosts := lipgloss.JoinVertical(lipgloss.Top, posts...)
	
	postList := m.postStyle.
	        Width(postsWidth).
	        MaxWidth(postsWidth).
		Render(renderedPosts)
	return postList
}


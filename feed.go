package main

import (
	"database/sql"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/bubbles/viewport"
)

type FindPostsFunc func(db *sql.DB, viewer SavedUser) ([]Post, error)

func getFeedView(renderer *lipgloss.Renderer, db *sql.DB, user SavedUser, find FindPostsFunc, name string) (Tab) {
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

	posts, err := find(db, user)
	if err != nil {
		log.Error(err)
	}
	timeline := getTimeline(renderer, db, posts, user)

	textInput := textarea.New()
	textInput.Placeholder = "Type a message..."

	textInput.CharLimit = 280
	textInput.SetWidth(30)
	textInput.SetHeight(3)
	textInput.FocusedStyle.CursorLine = lipgloss.NewStyle()
	textInput.ShowLineNumbers = false

	newViewport := viewport.New(20, 15)

	return Tab {
		Model: FeedModel{ 
			infoStyle: infoStyle, 
			quitStyle: quitStyle,
			postStyle: postStyle,
			headerStyle: headerStyle,
			numberStyle: numberStyle,
			db: db,
			renderer: renderer,
			posts: timeline,
			user: user,
			text: textInput,
			inputOpened: false,
			viewport: newViewport,
			find: find,
		},
		Name: name,
	}
}


type FeedModel struct {
	Name         string
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
	text         textarea.Model
	inputOpened  bool
	viewport     viewport.Model
	find         FindPostsFunc
}

func (m FeedModel) Init() tea.Cmd {
	return nil
}

func (m FeedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: 
		m.width = msg.Width
		m.posts.width = max(m.width, 20) - 2
		m.viewport.Width = m.posts.width
		m.viewport.Height = msg.Height - 5
		m.viewport.SetContent(m.posts.View())
		return m, nil
	case tea.KeyMsg:
		if m.text.Focused() {
			switch msg.String() {
			case "esc":
				m.text.Blur()
				m.inputOpened = false
				m.viewport.Height = m.viewport.Height + 4
				return m, nil
			case "enter":
				m.text.Blur()
				m.inputOpened = false
				m.viewport.Height = m.viewport.Height + 4
				text := m.text.Value()
				m.text.Reset()
				if (text == "") { 
					return m, nil
				}
				id, err := SavePost(m.db, m.user, text) 
				if err == nil {
					m.posts.Push(Post {
						id: id, 
						userId: m.user.id, 
						content: text, 
						username: m.user.username, 
						createdAt: time.Now(),
					})
					m.viewport.SetContent(m.posts.View())
					m.viewport.GotoTop()
				} else {
					log.Error(err)
				}

				return m, nil
			default:
				var cmd tea.Cmd
				m.text, cmd = m.text.Update(msg)
				return m, cmd
			}
		} else {
			switch msg.String() {
			case "p":
				m.inputOpened = true
				m.viewport.Height = m.viewport.Height - 4
				return m, m.text.Focus()
			case "r":
				posts, err := m.find(m.db, m.user)
				if err != nil {
					log.Error(err)
				}
				m.posts = getTimeline(m.renderer, m.db, posts, m.user)
				m.viewport.SetContent(m.posts.View())
				return m, nil
			case "k", "j":
				m.posts, m.viewport = UpdateTimeline(m.posts, m.viewport, msg)
				return m, nil
			case "a", "g":
				var cmd tea.Cmd
				m.posts, cmd = m.posts.Update(msg)
				return m, cmd
			case "l":
				var err error
				post := m.posts.posts[m.posts.currentPost]
				if post.liked {
					err = DeleteLike(m.db, m.user, post)
				} else {
					err = SaveLike(m.db, m.user, post)
				}
				
				if err == nil {
					m.posts.posts[m.posts.currentPost].liked = !post.liked
					if post.liked {
						m.posts.posts[m.posts.currentPost].likes = post.likes - 1 
					} else {
						m.posts.posts[m.posts.currentPost].likes = post.likes + 1 
					}
					m.viewport.SetContent(m.posts.View())
				}
				return m, nil
			}
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m FeedModel) View() string {
	postsWidth := max(m.width, 20)
	posts := make([]string, 0)
	if m.inputOpened {
		posts = append(posts, m.text.View() + "\n")
	}
	posts = append(posts, m.viewport.View())
	renderedPosts := lipgloss.JoinVertical(lipgloss.Top, posts...)
	
	postList := m.postStyle.
	        Width(postsWidth).
	        MaxWidth(postsWidth).
		Render(renderedPosts)
	return postList
}


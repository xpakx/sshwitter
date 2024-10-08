package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func getTimeline(renderer *lipgloss.Renderer, db *sql.DB, posts []Post, user SavedUser) (TimelineModel) {
	postStyle := renderer.NewStyle().
		BorderForeground(lipgloss.Color("8"))

	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))

	headerStyle := renderer.NewStyle().
		Foreground(lipgloss.Color("5"))
	numberStyle := quitStyle.
		Bold(true)

	textInput := textarea.New()
	textInput.Placeholder = "Type a message..."

	textInput.CharLimit = 280
	textInput.SetWidth(30)
	textInput.SetHeight(3)
	textInput.FocusedStyle.CursorLine = lipgloss.NewStyle()
	textInput.ShowLineNumbers = false

	return TimelineModel{ 
		quitStyle: quitStyle,
		postStyle: postStyle,
		headerStyle: headerStyle,
		numberStyle: numberStyle,
		db: db,
		posts: posts,
		user: user,
		currentPost: 0,
	}
}


type TimelineModel struct {
	quitStyle    lipgloss.Style
	postStyle    lipgloss.Style
	headerStyle  lipgloss.Style
	numberStyle  lipgloss.Style
	posts        []Post
	indices      []PostIndice
	user         SavedUser
	db           *sql.DB
	width        int
	currentPost  int
}

type PostIndice struct {
	start int
	len   int
}

func (m TimelineModel) Init() tea.Cmd {
	return nil
}


func (m TimelineModel) Update(msg tea.Msg) (TimelineModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: 
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down": 
			m.currentPost = min(m.currentPost + 1, len(m.posts)-1);
			return m, nil
		case "k", "up": 
			m.currentPost = max(m.currentPost - 1, 0);
			return m, nil
		case "a":
			if m.currentPost > len(m.posts) {
				return m, nil
			}
			username := m.posts[m.currentPost].username
			return m, openProfile(username)
		}
	}
	return m, nil
}

func (m *TimelineModel) View() string {
	posts := make([]string, 0)
	start := 0
	m.indices = make([]PostIndice, len(m.posts))
	for i, post := range m.posts {
		postView := m.postView(post, m.currentPost == i)
		posts = append(posts, postView)
		m.indices[i].start = start
		m.indices[i].len = strings.Count(postView, "\n") + 1
		start += m.indices[i].len
	}
	renderedPosts := lipgloss.JoinVertical(lipgloss.Top, posts...)
	
	postList := m.postStyle.
	        Width(m.width).
	        MaxWidth(m.width).
		Render(renderedPosts)
	return postList
}

func RelativeTime(t time.Time) string {
	now := time.Now()
	duration := now.Sub(t)

	switch {
	case duration < 30*time.Second:
		return "just now"
	case duration < time.Minute:
		return fmt.Sprintf("%.0fs ago", duration.Seconds())
	case duration < time.Hour:
		return fmt.Sprintf("%.0fm ago", duration.Minutes())
	case duration < 24*time.Hour:
		return fmt.Sprintf("%.0fh ago", duration.Hours())
	case now.Year() == t.Year():
		return t.Format("Jan 2")
	default:
		return t.Format("Jan 2, 2006")
	}
}

func (m TimelineModel) postView(post Post, current bool) string {
	doc := strings.Builder{}
	doc.WriteString(m.headerStyle.Render(post.username))
	doc.WriteString(m.quitStyle.Render(" · "))
	doc.WriteString(m.quitStyle.Render(RelativeTime(post.createdAt)))
	if current {
		doc.WriteString(m.quitStyle.Render(" !"))
	}
	doc.WriteString("\n")

	doc.WriteString(post.content)
	doc.WriteString("\n")
	if (post.liked) {
		doc.WriteString(m.quitStyle.Render("❤ "))
	}
	doc.WriteString(m.numberStyle.Render(strconv.Itoa(post.likes)))
	doc.WriteString(m.quitStyle.Render(" Likes"))
	doc.WriteString("  ")
	doc.WriteString(m.numberStyle.Render("0"))
	doc.WriteString(m.quitStyle.Render(" Replies"))
	doc.WriteString("\n")
	return doc.String()
}


func (m *TimelineModel) Push(post Post) {
	m.posts = append([]Post{post}, m.posts...)
}

func UpdateTimeline(posts TimelineModel, viewport viewport.Model, msg tea.KeyMsg) (TimelineModel, viewport.Model) {
	posts, _ = posts.Update(msg)
	if len(posts.posts) == 0 {
		return posts, viewport
	}
	start := viewport.YOffset
	end := start + viewport.Height
	curr := posts.indices[posts.currentPost]
	currStart := curr.start
	currEnd := curr.start + curr.len
	if currStart < start || currEnd > end {
		viewport.SetYOffset(currStart)
	}
	viewport.SetContent(posts.View())
	return posts, viewport
}

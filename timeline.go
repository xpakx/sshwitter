package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
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
	}
}


type TimelineModel struct {
	quitStyle    lipgloss.Style
	postStyle    lipgloss.Style
	headerStyle  lipgloss.Style
	numberStyle  lipgloss.Style
	posts        []Post
	user         SavedUser
	db           *sql.DB
	width        int
}

func (m TimelineModel) Init() tea.Cmd {
	return nil
}


func (m TimelineModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: 
		m.width = msg.Width
	}
	return m, nil
}

func (m TimelineModel) View() string {
	posts := make([]string, 0)
	for _, post := range m.posts {
		posts = append(posts, m.postView(post))
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

func (m TimelineModel) postView(post Post) string {
	doc := strings.Builder{}
	doc.WriteString(m.headerStyle.Render(post.username))
	doc.WriteString(m.quitStyle.Render(" · "))
	doc.WriteString(m.quitStyle.Render(RelativeTime(post.createdAt)))
	doc.WriteString("\n")

	doc.WriteString(post.content)
	doc.WriteString("\n")
	doc.WriteString(m.numberStyle.Render("0"))
	doc.WriteString(m.quitStyle.Render(" Likes"))
	doc.WriteString("  ")
	doc.WriteString(m.numberStyle.Render("0"))
	doc.WriteString(m.quitStyle.Render(" Replies"))
	doc.WriteString("\n")
	return doc.String()
}

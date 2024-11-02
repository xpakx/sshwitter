package main

import (
	"database/sql"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func getPostView(renderer *lipgloss.Renderer, db *sql.DB, postId int64, user SavedUser) (Tab) {
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

	post, postFound :=  GetPostById(db, postId, user.username)

	if (!postFound) {
		log.Info("No such post")
		// TODO: 404 page
	}
	isOwner := user.id == post.userId

	return Tab{
		Model: PostViewModel{ 
			infoStyle: infoStyle, 
			quitStyle: quitStyle,
			postStyle: postStyle,
			infoWidth: infoWidth,
			headerStyle: headerStyle,
			numberStyle: numberStyle,
			db: db,
			renderer: renderer,
			user: user,
			post: post,
			isOwner: isOwner,
		},
		Name: post.username,
	}
}


type PostViewModel struct {
	infoStyle    lipgloss.Style
	quitStyle    lipgloss.Style
	postStyle    lipgloss.Style
	headerStyle  lipgloss.Style
	numberStyle  lipgloss.Style
	post         Post
	user         SavedUser
	db           *sql.DB
	renderer     *lipgloss.Renderer
	width        int
	infoWidth    int
	isOwner      bool
}

func (m PostViewModel) Init() tea.Cmd {
	return nil
}


func (m PostViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: 
		m.width = msg.Width
		return m, nil
	}
	return m, nil
}

func (m PostViewModel) View() string {
	post := m.post
	doc := strings.Builder{}
	doc.WriteString(m.headerStyle.Render(post.username))
	doc.WriteString(m.quitStyle.Render(" · "))
	doc.WriteString(m.quitStyle.Render(RelativeTime(post.createdAt)))
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

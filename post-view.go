package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
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
		log.Infof("No post with id %d", postId)
		// TODO: 404 page
	}

	var parent Post
	hasParent := false

	if (post.parentId.Valid) {
		parentId := post.parentId.Int64
		parent, hasParent =  GetPostById(db, parentId, user.username)
		if (!hasParent) {
			log.Infof("No post with id %d", parentId)
		}
	}

	isOwner := user.id == post.userId

	tabName := fmt.Sprintf("%s %d", post.username, post.id)

	textInput := textarea.New()
	textInput.Placeholder = "Type a message..."

	textInput.CharLimit = 280
	textInput.SetWidth(30)
	textInput.SetHeight(3)
	textInput.FocusedStyle.CursorLine = lipgloss.NewStyle()
	textInput.ShowLineNumbers = false

	posts, err := FindReplies(db, post.id, user)
	if err != nil {
		log.Error(err)
	}
	timeline := getTimeline(renderer, db, posts, user)
	vp := viewport.New(20, 15)
	timeline.PushFront(post)
	if (hasParent) {
		timeline.PushFront(parent)
		timeline.Highlight(1)
	} else {
		timeline.Highlight(0)
	}

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
			textarea: textInput,
			inputOpened: false,
			parent: parent,
			hasParent: hasParent,
			posts: timeline,
			viewport: vp,
		},
		Name: tabName,
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
	textarea     textarea.Model
	inputOpened  bool
	hasParent    bool
	parent       Post
	posts        TimelineModel
	viewport     viewport.Model
}

func (m PostViewModel) Init() tea.Cmd {
	return nil
}


func (m PostViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: 
		m.width = msg.Width
		m.posts.width = max(m.width - (m.infoWidth + 1), 20) - 2
		m.viewport.Width = m.posts.width
		m.viewport.Height = msg.Height - 5
		m.viewport.SetContent(m.posts.View())
		return m, nil
	case tea.KeyMsg:
		if m.textarea.Focused() {
			switch msg.String() {
			case "esc":
				m.textarea.Blur()
				m.inputOpened = false
				m.viewport.Height = m.viewport.Height + 4
				return m, nil
			case "enter":
				m.textarea.Blur()
				m.inputOpened = false
				text := m.textarea.Value()
				m.viewport.Height = m.viewport.Height + 4
				m.textarea.Reset()
				if (text == "") { 
					return m, nil
				}
				ReplyToPost(m.db, m.user, m.post, text)
				return m, nil
			default:
				var cmd tea.Cmd
				m.textarea, cmd = m.textarea.Update(msg)
				return m, cmd
			}
		} else {
			switch msg.String() {
			case "p":
				m.inputOpened = true
				m.viewport.Height = m.viewport.Height - 4
				return m, m.textarea.Focus()
			}
		}

	}
	return m, nil
}

func (m PostViewModel) View() string {
	doc := strings.Builder{}
	if m.inputOpened {
		doc.WriteString(m.textarea.View())
		doc.WriteString("\n")
	}

	doc.WriteString(m.viewport.View())

	return doc.String()
}

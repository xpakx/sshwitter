package main

import (
	"database/sql"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/bubbles/textarea"
)

func getProfileView(renderer *lipgloss.Renderer, db *sql.DB, username string, user SavedUser) (ProfileViewModel) {
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

	owner, found :=  GetUserByUsername(db, username)
	if (!found) {
		log.Info("No such user")
		// TODO: 404 page
	}
	isOwner := user.id == owner.id

	info := getProfileInfo(renderer, db, owner)
	posts := make([]PostModel, 0)
	posts = append(posts, getPostView(renderer, owner))
	posts = append(posts, getPostView(renderer, owner))
	posts = append(posts, getPostView(renderer, owner))
	posts = append(posts, getPostView(renderer, owner))
	posts = append(posts, getPostView(renderer, owner))

	textInput := textarea.New()
	textInput.Placeholder = "Type a message..."

	textInput.CharLimit = 280
	textInput.SetWidth(30)
	textInput.SetHeight(3)
	textInput.FocusedStyle.CursorLine = lipgloss.NewStyle()
	textInput.ShowLineNumbers = false

	return ProfileViewModel{ 
		infoStyle: infoStyle, 
		quitStyle: quitStyle,
		postStyle: postStyle,
		infoWidth: infoWidth,
		owner: owner,
		db: db,
		info: info,
		posts: posts,
		user: user,
		isOwner: isOwner,
		text: textInput,
		inputOpened: false,
	}
}


type ProfileViewModel struct {
	infoStyle    lipgloss.Style
	quitStyle    lipgloss.Style
	postStyle    lipgloss.Style
	info         ProfileInfoModel
	posts        []PostModel
	owner         SavedUser
	user         SavedUser
	db           *sql.DB
	width        int
	infoWidth    int
	isOwner      bool
	text         textarea.Model
	inputOpened  bool
}

func (m ProfileViewModel) Init() tea.Cmd {
	return nil
}


func (m ProfileViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: 
		m.width = msg.Width
	case tea.KeyMsg:
		if m.text.Focused() {
			switch msg.String() {
			case "esc":
				m.text.Blur()
				m.inputOpened = false
				return m, nil
			default:
				var cmd tea.Cmd
				m.text, cmd = m.text.Update(msg)
				return m, cmd
			}
		} else {
			switch msg.String() {
			case "p":
				if m.isOwner {
					m.inputOpened = true
					return m, m.text.Focus()
				}
			}
		}
	}
	return m, nil
}

func (m ProfileViewModel) View() string {
	postsWidth := max(m.width - (m.infoWidth + 1), 20)
	info := m.infoStyle.
		Render(m.info.View())

	posts := make([]string, 0)
	if m.inputOpened {
		posts = append(posts, m.text.View() + "\n")
	}
	for _, post := range m.posts {
		posts = append(posts, post.View())
	}
	renderedPosts := lipgloss.JoinVertical(lipgloss.Top, posts...)
	
	postList := m.postStyle.
	        Width(postsWidth).
	        MaxWidth(postsWidth).
		Render(renderedPosts)
	doc := lipgloss.JoinHorizontal(lipgloss.Left,
		info,
		postList,
	)
	return doc
}

func getProfileInfo(renderer *lipgloss.Renderer, db *sql.DB, user SavedUser) (ProfileInfoModel) {
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
	headerStyle := renderer.NewStyle().
		Foreground(lipgloss.Color("5")).
		Bold(true)
	numberStyle := renderer.NewStyle().
		Bold(true)

	return ProfileInfoModel{ 
		txtStyle: txtStyle, 
		quitStyle: quitStyle,
		headerStyle: headerStyle,
		numberStyle: numberStyle,
		user: user,
		db: db,
	}
}


type ProfileInfoModel struct {
	txtStyle     lipgloss.Style
	quitStyle    lipgloss.Style
	headerStyle  lipgloss.Style
	numberStyle  lipgloss.Style
	user         SavedUser
	db           *sql.DB
}

func (m ProfileInfoModel) Init() tea.Cmd {
	return nil
}


func (m ProfileInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ProfileInfoModel) View() string {
	doc := strings.Builder{}
	username := m.headerStyle.Render(m.user.username)  
	doc.WriteString(username)
	doc.WriteString("\n")
	description := m.txtStyle.Render("description")  
	doc.WriteString(description)
	doc.WriteString("\n\n")
	doc.WriteString("📍 " )
	location := m.quitStyle.Render("City")  
	doc.WriteString(location)
	doc.WriteString("\n")
	doc.WriteString("🗓  " )
	joinDate := m.quitStyle.Render("Date")  
	doc.WriteString(joinDate)
	doc.WriteString("\n")
	doc.WriteString(m.numberStyle.Render("0"))
	doc.WriteString(m.quitStyle.Render(" Following"))
	doc.WriteString("\n")
	doc.WriteString(m.numberStyle.Render("0"))
	doc.WriteString(m.quitStyle.Render(" Followers"))
	return doc.String()
}







func getPostView(renderer *lipgloss.Renderer, user SavedUser) (PostModel) {
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))
	headerStyle := renderer.NewStyle().
		Foreground(lipgloss.Color("5"))
	numberStyle := quitStyle.
		Bold(true)

	return PostModel{ 
		txtStyle: txtStyle, 
		quitStyle: quitStyle,
		headerStyle: headerStyle,
		numberStyle: numberStyle,
		user: user,
	}
}

type PostModel struct {
	txtStyle     lipgloss.Style
	quitStyle    lipgloss.Style
	headerStyle  lipgloss.Style
	numberStyle  lipgloss.Style
	user         SavedUser
}

func (m PostModel) Init() tea.Cmd {
	return nil
}


func (m PostModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m PostModel) View() string {
	doc := strings.Builder{}
	doc.WriteString(m.headerStyle.Render(m.user.username))
	doc.WriteString(m.quitStyle.Render(" · "))
	doc.WriteString(m.quitStyle.Render("just now"))
	doc.WriteString("\n")

	doc.WriteString("post")
	doc.WriteString("\n")
	doc.WriteString(m.numberStyle.Render("0"))
	doc.WriteString(m.quitStyle.Render(" Likes"))
	doc.WriteString("  ")
	doc.WriteString(m.numberStyle.Render("0"))
	doc.WriteString(m.quitStyle.Render(" Replies"))
	doc.WriteString("\n")
	return doc.String()
}

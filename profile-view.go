package main

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
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

	headerStyle := renderer.NewStyle().
		Foreground(lipgloss.Color("5"))
	numberStyle := quitStyle.
		Bold(true)

	owner, found :=  GetUserByUsername(db, username)
	if (!found) {
		log.Info("No such user")
		// TODO: 404 page
	}
	isOwner := user.id == owner.id
	follows := false;
	if !isOwner {
		var err error;
		follows, err = CheckFollow(db, user, owner);
		if err != nil {
			log.Error("Error while checking following.")
		}
	} 


	posts, err := FindUserPosts(db, owner, user)
	if err != nil {
		log.Error(err)
	}
	timeline := getTimeline(renderer, db, posts, user)

	info := getProfileInfo(renderer, db, owner, follows)

	textInput := textarea.New()
	textInput.Placeholder = "Type a message..."

	textInput.CharLimit = 280
	textInput.SetWidth(30)
	textInput.SetHeight(3)
	textInput.FocusedStyle.CursorLine = lipgloss.NewStyle()
	textInput.ShowLineNumbers = false


	newViewport := viewport.New(20, 15)

	return ProfileViewModel{ 
		infoStyle: infoStyle, 
		quitStyle: quitStyle,
		postStyle: postStyle,
		infoWidth: infoWidth,
		headerStyle: headerStyle,
		numberStyle: numberStyle,
		owner: owner,
		db: db,
		renderer: renderer,
		info: info,
		posts: timeline,
		user: user,
		isOwner: isOwner,
		text: textInput,
		inputOpened: false,
		viewport: newViewport,
	}
}


type ProfileViewModel struct {
	infoStyle    lipgloss.Style
	quitStyle    lipgloss.Style
	postStyle    lipgloss.Style
	headerStyle  lipgloss.Style
	numberStyle  lipgloss.Style
	info         ProfileInfoModel
	posts        TimelineModel
	owner        SavedUser
	user         SavedUser
	db           *sql.DB
	renderer     *lipgloss.Renderer
	width        int
	infoWidth    int
	isOwner      bool
	text         textarea.Model
	inputOpened  bool
	viewport     viewport.Model
}

func (m ProfileViewModel) Init() tea.Cmd {
	return nil
}




func (m ProfileViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: 
		m.width = msg.Width
		m.posts.width = max(m.width - (m.infoWidth + 1), 20) - 2
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
				if (text == "") { // TODO
					return m, nil
				}
				id, err := SavePost(m.db, m.user, text) // TODO: extract to commad
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
				if m.isOwner {
					m.inputOpened = true
					m.viewport.Height = m.viewport.Height - 4
					return m, m.text.Focus()
				}
			case "r":
				posts, err := FindUserPosts(m.db, m.owner, m.user)
				if err != nil {
					log.Error(err)
				}
				m.posts = getTimeline(m.renderer, m.db, posts, m.user)
				m.viewport.SetContent(m.posts.View())
				return m, nil
			case "f":
				if m.isOwner {
					return m, nil
				}
				if !m.info.isFollowed {
					err := SaveFollow(m.db, m.user, m.owner)
					if err == nil {
						m.info.isFollowed = true
						m.owner.followers += 1
						m.info.user.followers += 1
					}
					return m, nil
				} else {
					err := DeleteFollow(m.db, m.user, m.owner)
					if err == nil {
						m.info.isFollowed = false
						m.owner.followers -= 1
						m.info.user.followers -= 1
					}
					return m, nil
				}
			}
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
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
	posts = append(posts, m.viewport.View())
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


func getProfileInfo(renderer *lipgloss.Renderer, db *sql.DB, user SavedUser, isFollowed bool) (ProfileInfoModel) {
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
		isFollowed: isFollowed,
		user: user,
		db: db,
	}
}


type ProfileInfoModel struct {
	txtStyle     lipgloss.Style
	quitStyle    lipgloss.Style
	headerStyle  lipgloss.Style
	numberStyle  lipgloss.Style
	isFollowed   bool
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
	if m.isFollowed {
		doc.WriteString("\n[Followed]")
	}
	doc.WriteString("\n\n")
	doc.WriteString("📍 " )
	location := m.quitStyle.Render("City")  
	doc.WriteString(location)
	doc.WriteString("\n")
	doc.WriteString("🗓  " )
	joinDate := m.quitStyle.Render("Date")  
	doc.WriteString(joinDate)
	doc.WriteString("\n")
	doc.WriteString(m.numberStyle.Render(strconv.Itoa(m.user.followed)))
	doc.WriteString(m.quitStyle.Render(" Following"))
	doc.WriteString("\n")
	doc.WriteString(m.numberStyle.Render(strconv.Itoa(m.user.followers)))
	doc.WriteString(m.quitStyle.Render(" Followers"))
	return doc.String()
}

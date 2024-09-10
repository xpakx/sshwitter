package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
)


// TODO: move to yaml/env
const (
	host          = "localhost"
	port          = "23230"
	validPassword = "password"
	dev           = true
)

// TODO: move to db
var users = map[string]string{
}

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithPublicKeyAuth(GetPublicKeyAuth),

		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)

	if (dev) { s.SetOption(wish.WithPasswordAuth(GetPasswordAuth)) }


	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}


func GetPasswordAuth(_ ssh.Context, password string) bool {
	if (password == validPassword) {
		log.Info("Successful authentication")
	} else {
		log.Info("Authentication failed")
	}
	return password == validPassword
}

func GetPublicKeyAuth(_ ssh.Context, key ssh.PublicKey) bool {
	log.Info("public-key")
	for _, pubkey := range users {
		parsed, _, _, _, _ := ssh.ParseAuthorizedKey(
			[]byte(pubkey),
		)
		if ssh.KeysEqual(key, parsed) {
			return true
		}
	}
	return false
}


func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	renderer := bubbletea.MakeRenderer(s)
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))

	m := model{ name: "sshwitter", txtStyle: txtStyle, quitStyle: quitStyle }
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

type model struct {
	name      string
	txtStyle  lipgloss.Style
	quitStyle  lipgloss.Style
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("Authorized!")
	return m.txtStyle.Render(s) + "\n\n" + m.quitStyle.Render("Press 'q' to quit\n")
}

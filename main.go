package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
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

		wish.WithPublicKeyAuth(func(_ ssh.Context, key ssh.PublicKey) bool {
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
		}),

		wish.WithMiddleware(
			logging.Middleware(),
			func(next ssh.Handler) ssh.Handler {
				return func(sess ssh.Session) {
					wish.Println(sess, "Authorized!")
				}
			},
		),
	)


	if (dev) {
		s.SetOption(
			wish.WithPasswordAuth(func(_ ssh.Context, password string) bool {
				if (password == validPassword) {
					log.Info("Successful authentication")
				} else {
					log.Info("Authentication failed")
				}
				return password == validPassword
			}),
		)
	}


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

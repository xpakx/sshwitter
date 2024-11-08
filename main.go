package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"

	gossh "golang.org/x/crypto/ssh"

	"database/sql"
	_ "github.com/lib/pq"
)

// TODO: move to yaml/env
const (
	defaultHost          = "localhost"
	port          = "23230"
	validPassword = "password"
	dev           = false
)

func GetEnvOrDefault(env string, def string) string {
	result, defined := os.LookupEnv(env)
	if defined {
		return result
	} else {
		return def
	}
}

func main() {
	host := GetEnvOrDefault("HOST", defaultHost)
	dbUrl := GetEnvOrDefault("DB",  "postgresql://root:password@localhost:5432/sshwitter?sslmode=disable")

	db, dbErr := openDB(dbUrl)
	if dbErr != nil {
		log.Error(dbErr.Error())
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		log.Info("Database connection closed")
		if err != nil {
			log.Error(err.Error())
		}
	}(db)

	CreateUserTable(db)
	CreatePostTable(db)
	CreateFollowTable(db)
	CreateFollowFunction(db)
	CreateUnfollowFunction(db)
	CreateLikeTable(db)
	CreateLikeFunction(db)
	CreateUnlikeFunction(db)
	CreateReplyFunction(db)

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithPublicKeyAuth(makeGetPublicKeyAuth(db)),

		wish.WithMiddleware(
			bubbletea.Middleware(makeTeaHandler(db)),
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


func GetPasswordAuth(context ssh.Context, password string) bool {
	if (password == validPassword) {
		log.Info("Successful authentication")
	} else {
		log.Info("Authentication failed")
	}
	context.SetValue("guest", true);
	context.SetValue("verified", false);
	context.SetValue("publicKey", "");
	return password == validPassword
}

func ConvertKey(k ssh.PublicKey) string {
	return strings.TrimSpace(string(gossh.MarshalAuthorizedKey(k)[:]))
}

func GetPublicKeyAuth(context ssh.Context, db *sql.DB, key ssh.PublicKey) bool {
	username := strings.Split(context.User(), ":")[0];
	log.Infof("New connection with username: %s", username)
	log.Info("Trying public key")

	if savedUser, found := GetUserByUsername(db, username); found {
		parsed, _, _, _, _ := ssh.ParseAuthorizedKey(
			[]byte(savedUser.key),
		)
		
		if ssh.KeysEqual(key, parsed) {
			context.SetValue("guest", false);
			context.SetValue("verified", savedUser.verified);
			context.SetValue("user", savedUser);
			return true
		}
	} 
	log.Info("Public key not found")
	context.SetValue("guest", true);
	context.SetValue("publicKey", ConvertKey(key));
	context.SetValue("verified", false);
	return true
}


func teaHandler(s ssh.Session, db *sql.DB) (tea.Model, []tea.ProgramOption) {
	username := strings.Split(s.Context().User(), ":")[0];
	guest := s.Context().Value("guest").(bool)
	verified := s.Context().Value("verified").(bool)
	log.Info("New tea handler", "user", username, "guest", guest, "verified", verified)

	renderer := bubbletea.MakeRenderer(s)

	var model tea.Model

	if (!guest && verified) {
		user := s.Context().Value("user").(SavedUser)
		model =  getBoardModel(renderer, db, user)
	} else if (!guest && !verified) {
		model = getUnverifiedModel(renderer, username)
	} else {
		publicKey := s.Context().Value("publicKey").(string)
		model = getRegisterModel(renderer, db, username, publicKey)
	}
	return model, []tea.ProgramOption{tea.WithAltScreen()}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func makeTeaHandler(db *sql.DB) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
    return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
        return teaHandler(s, db)
    }
}

func makeGetPublicKeyAuth(db *sql.DB) func(context ssh.Context, key ssh.PublicKey) bool {
	return func(context ssh.Context, key ssh.PublicKey) bool {
		return GetPublicKeyAuth(context, db, key)
	}
}

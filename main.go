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
)

// TODO: move to yaml/env
const (
	host          = "localhost"
	port          = "23230"
	validPassword = "password"
	dev           = false
)


type SavedUser struct {
	key             string
	verified        bool
	administrator   bool
	email           string
	username        string
}

// TODO: move to db
var users = map[string]SavedUser{
	"test" : { 
		key: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDQ2EOmuzeJo6LN4jg7Vuvf0BSjrqDMSguc3Z6i0zkVURz5Bb61BZlG7PpyqOQ1aeISNfoMSpVwx1fQytIvQBfpF6OX+XMI/wzgOEvbPNeQ0RHsPJjq0x6wMLNEWpPl5f07pDgBlWB8IkBTKvSZQje/WsEwDvUnFRrWcC8PHs2H/WRpm+wagg9T5N6jDqlC711DJEWIyKwl744QHK4NBnyXHfK+0pW/JfhEelyQ+bTVfWNDu9V5uZI69hiKZNs4UANhAoUEhhZIy60ZHho6Zn8JkZkjORMwGi/hi8lUaIDYXXcqKGqKQdU2HU5NgpWVO3/w7KRQceegDiMO5Aa/yMEtdVi0B2NmUGVTZcCEwkqWbACqG5r23AmgrMX/Hh8L/9Z1nFwnxCY2bUd29DQI1q7GzTwYIxNi9y7/8H5+gmU6Yn3Wm5mUjpxWLF9QbU0fOFNZ/WO1h3rRYCwoouJ4ixWuCM6BLcBuuEutx24mjBaO3x0p68XJ8rxMuvS/n9TwTywPfeDS5Yft1hHovyRt1vSAHLxd8eSP65vJHJwsYAL8psGbm68CyYnzf8D4CPSJh4DSCQRzNnfFjYozX9QuXAhPtJkjPI7w6mJyPmjUaDB+sOkolIqIdF0jBXuaB/Hv/03H3ul5+SqpB0s37Wh0rwI2ORX0Ct45pYjj78WtAkukEQ==",
		verified: true,
		administrator: true,
		email: "",
		username: "test",
	},
	"test2" : {verified: false, username: "test2"},
	"test3" : {verified: false, username: "test3"},
	"test4" : {verified: false, username: "test4"},
	"test5" : {verified: false, username: "test5"},
	"test6" : {verified: false, username: "test6"},
	"test7" : {verified: false, username: "test7"},
	"test8" : {verified: false, username: "test8"},
	"test9" : {verified: false, username: "test9"},
	"test10" : {verified: false, username: "test10"},
	"test11" : {verified: false, username: "test11"},
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

func GetPublicKeyAuth(context ssh.Context, key ssh.PublicKey) bool {
	username := strings.Split(context.User(), ":")[0];
	log.Infof("New connection with username: %s", username)
	log.Info("Trying public key")

	if savedUser, found :=  users[username]; found {
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


func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	username := strings.Split(s.Context().User(), ":")[0];
	guest := s.Context().Value("guest").(bool)
	verified := s.Context().Value("verified").(bool)
	log.Info("New tea handler", "user", username, "guest", guest, "verified", verified)

	renderer := bubbletea.MakeRenderer(s)

	var model tea.Model

	if (!guest && verified) {
		user := s.Context().Value("user").(SavedUser)
		model =  getBoardModel(renderer, user)
	} else if (!guest && !verified) {
		model = getUnverifiedModel(renderer, username)
	} else {
		publicKey := s.Context().Value("publicKey").(string)
		model = getRegisterModel(renderer, username, publicKey)
	}
	return model, []tea.ProgramOption{tea.WithAltScreen()}
}

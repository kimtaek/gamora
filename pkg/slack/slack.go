package slack

import (
	"encoding/json"
	"github.com/caarlos0/env"
	"net/http"
	"net/url"
)

// Configure config for slack
type Configure struct {
	URL string `env:"SLACK_URL" envDefault:"https://hooks.slack.com/services/zzz/xxx/yyy"`
}

// Message defined message for slack
type Message struct {
	Text     string `json:"text"`
	Username string `json:"username"`
	Channel  string `json:"channel"`
}

// Config global defined slack config
var Config Configure

// Setup init slack config
func Setup() {
	_ = env.Parse(&Config)
}

// SendMessage send slack message
func SendMessage(s Message) {
	if s.Text == "" {
		return
	}
	if s.Channel == "" {
		s.Channel = ExceptionChannel()
	}
	p, _ := json.Marshal(s)
	r, _ := http.PostForm(
		Config.URL,
		url.Values{
			"payload": {string(p)},
		},
	)
	defer r.Body.Close()
}

// ExceptionChannel defined default exception channel
func ExceptionChannel() string {
	return "#exception"
}

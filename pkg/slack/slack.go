package slack

import (
	"encoding/json"
	"github.com/caarlos0/env"
	"net/http"
	"net/url"
)

type config struct {
	URL string `env:"SLACK_URL" envDefault:"https://hooks.slack.com/services/zzz/xxx/yyy"`
}

type Message struct {
	Text     string `json:"text"`
	Username string `json:"username"`
	Channel  string `json:"channel"`
}

var c config

func Setup() {
	_ = env.Parse(&c)
}

// SendSlackMessage
func SendMessage(s Message) {
	if s.Text == "" {
		return
	}
	if s.Channel == "" {
		s.Channel = ExceptionChannel()
	}
	p, _ := json.Marshal(s)
	r, _ := http.PostForm(
		c.URL,
		url.Values{
			"payload": {string(p)},
		},
	)
	defer r.Body.Close()
}

func ExceptionChannel() string {
	return "#exception"
}

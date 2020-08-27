package app

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/kimtaek/gamora/pkg/helper"
)

// Configure config for application
type Configure struct {
	Mode     string `env:"APP_MODE" envDefault:"debug"`
	AppHost  string `env:"APP_HOST" envDefault:"localhost"`
	TimeZone string `env:"TIMEZONE" envDefault:"Asia/Seoul"`
}

// Config global defined application config
var Config Configure

// Setup init application config
func Setup() {
	_ = env.Parse(&Config)
	helper.Info(fmt.Sprintf("running %s mode!", Config.Mode))
}

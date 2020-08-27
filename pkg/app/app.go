package app

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/kimtaek/gamora/pkg/helper"
)

type Configure struct {
	Version  int    `env:"VERSION" envDefault:"1"`
	Mode     string `env:"APP_MODE" envDefault:"debug"`
	AppHost  string `env:"APP_HOST" envDefault:"localhost"`
	TimeZone string `env:"TIMEZONE" envDefault:"Asia/Seoul"`
}

var Config Configure

func Setup() {
	_ = env.Parse(&Config)
	helper.Info(fmt.Sprintf("running %s mode!", Config.Mode))
}

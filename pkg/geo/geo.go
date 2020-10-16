package geo

import (
	"github.com/caarlos0/env"
	"github.com/oschwald/geoip2-golang"
)

// Configure config for geo
type Configure struct {
	GeoDatabasePath string `env:"GEO_DATABASE_PATH" envDefault:"./GeoIP2-City.mmdb"`
}

// Config global defined geo config
var Config Configure

// Database global defined geo database
var Database *geoip2.Reader

// Setup init geo config
func Setup() {
	_ = env.Parse(&Config)
	Database, _ = geoip2.Open(Config.GeoDatabasePath)
}

func CloseGeo() {
	defer Database.Close()
}

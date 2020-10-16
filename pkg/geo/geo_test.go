package geo

import (
	"github.com/kimtaek/gamora/pkg/helper"
	"net"
	"testing"
)

func TestCity(t *testing.T) {
	Setup()
	record, _ := Database.City(net.ParseIP("211.202.26.28"))
	helper.Info(record)
	defer Database.Close()
}

package tests

import (
	"testing"

	"gorm.io/gorm"

	"github.com/opentreehole/backend/cmd/wire"
)

var DB *gorm.DB

func init() {
	server, _, err := wire.NewApp()
	if err != nil {
		panic(err)
	}

	RegisterApp(server.GetFiberApp())
	DB = server.GetDB()
}

func TestAuth(t *testing.T) {
	DefaultTester.Get(t, "/docs/index.html", 200, RequestConfig{})
}

package auth

import (
	"testing"

	"github.com/opentreehole/backend/internal/auth"
	"github.com/opentreehole/backend/internal/pkg/server"
	"github.com/opentreehole/backend/tests"
)

func TestAuth(t *testing.T) {
	tests.RegisterApp(server.Init(auth.Config))

	tests.DefaultTester.Get(t, "/docs/index.html", 200, tests.RequestConfig{})
}

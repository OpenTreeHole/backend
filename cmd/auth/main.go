// @title Open Tree Hole Auth
// @version 3.0.0
// @description Next Generation of Auth microservice integrated with kong for registration and issuing tokens

// @contact.name Maintainer Chen Ke
// @contact.url https://danxi.fduhole.com/about
// @contact.email dev@fduhole.com

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000

package main

import (
	_ "time/tzdata"

	"github.com/opentreehole/backend/internal/auth"
	_ "github.com/opentreehole/backend/internal/auth/docs"
	"github.com/opentreehole/backend/internal/pkg/server"
)

func main() {
	server.Run(auth.Config)
}

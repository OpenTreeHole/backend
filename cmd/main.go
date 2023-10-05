// @title OpenTreeHole Backend
// @version 3.0.0
// @description Next Generation of OpenTreeHole Backend

// @contact.name Maintainer Chen Ke
// @contact.url https://danxi.fduhole.com/about
// @contact.email dev@fduhole.com

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api
// @schemes http https

package main

import (
	_ "time/tzdata"

	"github.com/spf13/cobra"

	"github.com/opentreehole/backend/cmd/migrate"
	"github.com/opentreehole/backend/cmd/wire"
	_ "github.com/opentreehole/backend/internal/docs"
)

var rootCmd = &cobra.Command{
	Use: "opentreehole_backend",
	Run: func(cmd *cobra.Command, args []string) {
		server, cleanup, err := wire.NewApp()
		if err != nil {
			panic(err)
		}

		defer cleanup()
		server.Run()
	},
}

func init() {
	rootCmd.AddCommand(migrate.Cmd)
}

//go:generate wire gen ./wire
func main() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

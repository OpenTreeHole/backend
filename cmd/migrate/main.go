//! 迁移旧版本蛋壳到新版本
//! `review`.`history` -> `review_history`

package migrate

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "migrate",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("no migration specified")
		}

		switch args[0] {
		case "danke_v3":
			DankeV3()
		default:
			return fmt.Errorf("unknown migration %s", args[0])
		}
		return nil
	},
}

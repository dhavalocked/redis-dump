package terminal

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

// Execute initialize the command execution
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var pattern string

var copyCmd = &cobra.Command{
	Use:  "copy sourceRedis targetRedis",
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		copier(cmd, args, pattern)
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	// pattern to search for
	copyCmd.Flags().StringVar(&pattern, "queryString", "*", "")
}

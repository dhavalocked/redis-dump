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
var scanCount, report, exportRoutines, pushRoutines int
var override bool

var copyCmd = &cobra.Command{
	Use:  "copy sourceRedis targetRedis",
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		copier(cmd, args, pattern, scanCount, exportRoutines, pushRoutines, override)
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	// pattern to search for
	copyCmd.Flags().StringVar(&pattern, "queryString", "*", "")
	// how much should i scan at once?
	copyCmd.Flags().IntVar(&scanCount, "scanLimit", 1000, "")
	// total threads i want to dump ?
	copyCmd.Flags().IntVar(&exportRoutines, "dumpThreads", 100, "")
	// total threads i want to restore ?
	copyCmd.Flags().IntVar(&pushRoutines, "restoreThreads", 100, "")
	// Do i want to override the keys if it exist in destination?
	copyCmd.Flags().BoolVar(&override, "overrideKey", false, "")
}

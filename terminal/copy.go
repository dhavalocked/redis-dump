package terminal

import (
	"fmt"
	"log"
	"redis-dump/pkg/logger"
	"redis-dump/pkg/restore"
	"redis-dump/pkg/scanner"
	"sync"
	"time"

	"github.com/mediocregopher/radix/v3"
	"github.com/spf13/cobra"
)

func copier(cmd *cobra.Command, args []string, pattern string) {
	fmt.Println("======== STARATING =========")

	scanCount := 1000
	exportRoutines := 300
	pushRoutines := 300

	clientSource, err := radix.DefaultClientFunc("tcp", args[0])
	if err != nil {
		log.Fatal(err)
	}

	clientTarget, err := radix.DefaultClientFunc("tcp", args[1])
	if err != nil {
		log.Fatal(err)
	}

	statusReporter := logger.NewLogger()

	redisScanner := scanner.NewScanner(
		clientSource,
		scanner.RedisScannerOpts{
			Pattern:          pattern,
			ScanCount:        scanCount,
			PullRoutineCount: exportRoutines,
		},
		statusReporter,
	)

	redisPusher := restore.NewRedisPusher(clientTarget, redisScanner.GetDumpChannel(), statusReporter)

	waitingGroup := new(sync.WaitGroup)

	// Log it every 5 seconds
	statusReporter.Start(time.Second * time.Duration(5))
	redisPusher.Start(waitingGroup, pushRoutines)
	redisScanner.Start()

	waitingGroup.Wait()
	statusReporter.Stop()
	statusReporter.Report()

	fmt.Println("======== DONE =========")
}

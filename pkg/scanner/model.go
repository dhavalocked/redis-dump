package scanner

import (
	"redis-dump/pkg/logger"

	"github.com/mediocregopher/radix/v3"
)

// KeyDump ...
type KeyDump struct {
	Key   string
	Value string
	Ttl   int
}

// RedisScannerOpts ...
type RedisScannerOpts struct {
	Pattern          string
	ScanCount        int
	PullRoutineCount int
}

// RedisScanner ...
type RedisScanner struct {
	client      radix.Client
	options     RedisScannerOpts
	reporter    *logger.Logger
	keyChannel  chan string
	dumpChannel chan KeyDump
}

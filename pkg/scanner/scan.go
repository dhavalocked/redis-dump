package scanner

import (
	"log"
	"redis-dump/pkg/logger"
	"sync"

	"github.com/mediocregopher/radix/v3"
)

func NewScanner(client radix.Client, options RedisScannerOpts, reporter *logger.Logger) *RedisScanner {
	return &RedisScanner{
		client:      client,
		options:     options,
		reporter:    reporter,
		dumpChannel: make(chan KeyDump),
		keyChannel:  make(chan string),
	}
}

func (s *RedisScanner) Start() {
	wgPull := new(sync.WaitGroup)
	wgPull.Add(s.options.PullRoutineCount)

	go s.scanRoutine()
	for i := 0; i < s.options.PullRoutineCount; i++ {
		go s.exportRoutine(wgPull)
	}

	wgPull.Wait()
	close(s.dumpChannel)
}

func (s *RedisScanner) GetDumpChannel() <-chan KeyDump {
	return s.dumpChannel
}

func (s *RedisScanner) scanRoutine() {
	var key string
	scanOpts := radix.ScanOpts{
		Command: "SCAN",
		Count:   s.options.ScanCount,
	}

	if s.options.Pattern != "*" {
		scanOpts.Pattern = s.options.Pattern
	}

	radixScanner := radix.NewScanner(s.client, scanOpts)
	for radixScanner.Next(&key) {
		s.keyChannel <- key
	}

	close(s.keyChannel)
}

func (s *RedisScanner) exportRoutine(wg *sync.WaitGroup) {
	for key := range s.keyChannel {
		var value string
		var ttl int

		p := radix.Pipeline(
			radix.Cmd(&ttl, "PTTL", key),
			radix.Cmd(&value, "DUMP", key),
		)

		if err := s.client.Do(p); err != nil {
			log.Fatal(err)
		}

		if ttl < 0 {
			ttl = 0
		}

		s.dumpChannel <- KeyDump{
			Key:   key,
			Ttl:   ttl,
			Value: value,
		}
	}

	wg.Done()
}

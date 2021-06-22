package restore

import (
	"fmt"
	"redis-dump/pkg/logger"
	"redis-dump/pkg/scanner"
	"sync"
	"time"

	"github.com/mediocregopher/radix/v3"
)

func NewRedisPusher(client radix.Client, dumpChannel <-chan scanner.KeyDump, reporter *logger.Logger) *RedisPusher {
	return &RedisPusher{
		client:      client,
		reporter:    reporter,
		dumpChannel: dumpChannel,
	}
}

type RedisPusher struct {
	client      radix.Client
	reporter    *logger.Logger
	dumpChannel <-chan scanner.KeyDump
}

func (p *RedisPusher) Start(wg *sync.WaitGroup, number int, overrideKey bool, timeout int) {
	wg.Add(number)
	for i := 0; i < number; i++ {
		go p.pushRoutine(overrideKey, timeout, wg)
	}

}

func (p *RedisPusher) pushRoutine(overrideKey bool, timeout int, wg *sync.WaitGroup) {
	for dump := range p.dumpChannel {
		p.reporter.AddCounter(1)

		if overrideKey {
			err := p.client.Do(radix.FlatCmd(nil, "RESTORE", dump.Key, dump.Ttl, dump.Value, "REPLACE"))
			if err != nil {
				fmt.Println("Got error while Restoring key %s on destination", dump.Key)
				fmt.Println(err)
			}
		} else {
			err := p.client.Do(radix.FlatCmd(nil, "RESTORE", dump.Key, dump.Ttl, dump.Value))
			if err != nil {
				if err.Error() == "BUSYKEY Target key name already exists" {
					fmt.Println("%s Key already exist.. ignoring ", dump.Key)
				} else {
					fmt.Println("Got error while Restoring key %s on destination", dump.Key)
					fmt.Println(err)
				}
			}
		}

		// add timeout of 10ms
		time.Sleep(20)
	}

	wg.Done()
}

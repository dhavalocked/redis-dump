package logger

import (
	"log"
	"sync/atomic"
	"time"
)

// Logger ...
type Logger struct {
	processedCount uint64
	start          time.Time
	doneChannel    chan bool
}

// NewLogger instance
func NewLogger() *Logger {
	return &Logger{
		doneChannel: make(chan bool),
	}
}

// Start ...
func (r *Logger) Start(reportPeriod time.Duration) {
	atomic.StoreUint64(&r.processedCount, 0)
	r.start = time.Now()
	go r.reportingRoutine(reportPeriod)
}

// Stop ...
func (r *Logger) Stop() {
	r.doneChannel <- true
}

// AddCounter ...
func (r *Logger) AddCounter(delta uint64) {
	atomic.AddUint64(&r.processedCount, delta)
}

// Report ...
func (r *Logger) Report() {
	log.Printf(
		"Processed & Dumped: %d  =====>  Total time taken %s \n",
		atomic.LoadUint64(&r.processedCount),
		time.Since(r.start),
	)
}

// reportingRoutine ...
func (r *Logger) reportingRoutine(reportPeriod time.Duration) {
	timer := time.NewTicker(reportPeriod)
	for {
		select {
		case <-timer.C:
			r.Report()
		case <-r.doneChannel:
			timer.Stop()
			break
		}
	}
}

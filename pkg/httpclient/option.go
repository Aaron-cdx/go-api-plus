package httpclient

import (
	"go-api-plus/pkg/trace"
	"sync"
	"time"

	"go.uber.org/zap"
)

var (
	// cache used for local cache
	cache = &sync.Pool{New: func() interface{} {
		return &option{header: make(map[string][]string)}
	}}
)

// Mock define the mock api data
type Mock func() (body []byte)

// Option http request custom settings
type Option func(*option)

type option struct {
	ttl         time.Duration
	header      map[string][]string
	trace       *trace.Trace
	dialog      *trace.Dialog
	logger      *zap.Logger
	retryTimes  int
	retryDelay  time.Duration
	retryVerify RetryVerify
	alarmTitle  string
	alarmObject AlarmObject
	alarmVerify AlarmVerify
	mock        Mock
}

func (o *option) reset() {
	o.ttl = 0
	o.header = make(map[string][]string)
	o.trace = nil
	o.dialog = nil
	o.logger = nil
	o.retryTimes = 0
	o.retryDelay = 0
	o.retryVerify = nil
	o.alarmTitle = ""
	o.alarmObject = nil
	o.alarmVerify = nil
	o.mock = nil
}

func getOption() *option {
	return cache.Get().(*option)
}

func releaseOption(opt *option) {
	opt.reset()
	cache.Put(opt)
}

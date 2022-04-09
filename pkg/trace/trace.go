package trace

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"sync"
)

const Header = "TRACE_ID"

var _ T = (*Trace)(nil)

type T interface {
	i()
	ID() string
	WithRequest(req *Request) *Trace
	WithResponse(resp *Response) *Trace
	AppendDialog(dialog *Dialog) *Trace
	AppendDebug(debug *Debug) *Trace
	AppendSQL(sql *SQL) *Trace
	AppendRedis(redis *Redis) *Trace
}

// Trace used for trace ann record the request info
type Trace struct {
	mux                sync.Mutex
	Identifier         string    `json:"trace_id"`
	Request            *Request  `json:"request"`
	Response           *Response `json:"response"`
	ThirdPartyRequests []*Dialog `json:"third_party_requests"`
	Debugs             []*Debug  `json:"debugs"`
	SQLs               []*SQL    `json:"sqls"`
	Redis              []*Redis  `json:"redis"`
	Success            bool      `json:"success"`
	CostSeconds        float64   `json:"cost_seconds"`
}

// Request request info
type Request struct {
	TTL        string      `json:"ttl"`         // request timeout time
	Method     string      `json:"method"`      // request method
	DecodedURL string      `json:"decoded_url"` // request url
	Header     interface{} `json:"header"`      // request header info
	Body       interface{} `json:"body"`        // request body info
}

// Response response info
type Response struct {
	Header          interface{} `json:"header"`                      // header info
	Body            interface{} `json:"body"`                        // body info
	BusinessCode    int         `json:"business_code,omitempty"`     // business code
	BusinessCodeMsg string      `json:"business_code_msg,omitempty"` // business tips info
	HttpCode        int         `json:"http_code"`                   // HTTP status code
	HttpCodeMsg     string      `json:"http_code_msg"`               // HTTP status code msg
	CostSeconds     float64     `json:"cost_seconds"`                // execute cost time(unit second)
}

// New used for get Trace instance to trace request process info
func New(id string) *Trace {
	if id == "" {
		buf := make([]byte, 10)
		if _, err := io.ReadFull(rand.Reader, buf); err != nil {
			panic("error reading from random source: " + err.Error())
		}
		id = hex.EncodeToString(buf)
	}
	return &Trace{Identifier: id}
}

func (t *Trace) i() {}

// ID identifier
func (t *Trace) ID() string {
	return t.Identifier
}

// WithRequest used for setting request
func (t *Trace) WithRequest(req *Request) *Trace {
	t.Request = req
	return t
}

// WithResponse used for setting response
func (t *Trace) WithResponse(resp *Response) *Trace {
	t.Response = resp
	return t
}

// AppendDialog used for append dialog info
func (t *Trace) AppendDialog(dialog *Dialog) *Trace {
	if dialog == nil {
		return t
	}
	t.mux.Lock()
	defer t.mux.Unlock()

	t.ThirdPartyRequests = append(t.ThirdPartyRequests, dialog)
	return t
}

// AppendDebug used for append debug info
func (t *Trace) AppendDebug(debug *Debug) *Trace {
	if debug == nil {
		return t
	}

	t.mux.Lock()
	defer t.mux.Unlock()

	t.Debugs = append(t.Debugs, debug)
	return t
}

// AppendSQL used for append sql process info
func (t *Trace) AppendSQL(sql *SQL) *Trace {
	if sql == nil {
		return t
	}

	t.mux.Lock()
	defer t.mux.Unlock()

	t.SQLs = append(t.SQLs, sql)
	return t
}

// AppendRedis used for append redis process info
func (t *Trace) AppendRedis(redis *Redis) *Trace {
	if redis == nil {
		return t
	}

	t.mux.Lock()
	defer t.mux.Unlock()

	t.Redis = append(t.Redis, redis)
	return t
}

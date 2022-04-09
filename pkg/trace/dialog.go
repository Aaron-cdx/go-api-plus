package trace

import "sync"

type D interface {
	i()
	AppendResponse(resp *Response)
}

// Dialog used for internal call others api dialog info, if failed, will retry
// so response will have many.
type Dialog struct {
	mux         sync.Mutex
	Request     *Request    `json:"request"`      // request info
	Responses   []*Response `json:"responses"`    // response info
	Success     bool        `json:"success"`      // whether success, false or true
	CostSeconds float64     `json:"cost_seconds"` // execute time(unit second)
}

func (d *Dialog) i() {}

// AppendResponse append response info
func (d *Dialog) AppendResponse(resp *Response) {
	if resp == nil {
		return
	}
	d.mux.Lock()
	d.Responses = append(d.Responses, resp)
	d.mux.Unlock()
}

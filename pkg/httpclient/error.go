package httpclient

var _ ReplyErr = (*replyErr)(nil)

// ReplyErr used for http error reply
type ReplyErr interface {
	error
	StatusCode() int
	Body() []byte
}

type replyErr struct {
	err        error
	statusCode int
	body       []byte
}

// Error return err info
func (r *replyErr) Error() string {
	return r.err.Error()
}

// StatusCode return status code info
func (r *replyErr) StatusCode() int {
	return r.statusCode
}

// Body return response body info
func (r *replyErr) Body() []byte {
	return r.body
}

func newReplyErr(statusCode int, body []byte, err error) ReplyErr {
	return &replyErr{
		err:        err,
		statusCode: statusCode,
		body:       body,
	}
}

// ConvertToReplyErr convert error to ReplyErr
func ConvertToReplyErr(err error) (ReplyErr, bool) {
	if err == nil {
		return nil, false
	}
	e, ok := err.(ReplyErr)
	return e, ok
}

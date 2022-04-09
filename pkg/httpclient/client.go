package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"go-api-plus/pkg/errors"
	"go-api-plus/pkg/trace"
	"net/url"
	"time"
)

var ErrorUrlRequired = errors.New("url required")

const (
	// DefaultTTL once a http request wait time
	DefaultTTL = time.Minute
	Json       = "application/json; charset=utf-8"
	Form       = "application/x-www-form-urlencoded; charset=utf-8"
)

// Get send get request
func Get(URL string, form url.Values, options ...Option) (body []byte, err error) {
	return
}

// withoutBody without or with body request
func withoutBody(method, URL string, form url.Values, options ...Option) (body []byte, err error) {
	if URL == "" {
		return nil, ErrorUrlRequired
	}

	if len(form) > 0 {
		if URL, err = addFormValuesIntoURL(URL, form); err != nil {
			return
		}
	}

	_ts := time.Now()
	opt := getOption()
	defer func() {
		if opt.trace != nil {
			opt.dialog.Success = err == nil
			opt.dialog.CostSeconds = time.Since(_ts).Seconds()
			opt.trace.AppendDialog(opt.dialog)
		}

		releaseOption(opt)
	}()

	for _, f := range options {
		f(opt)
	}

	opt.header["Content-Type"] = []string{
		"application/x-www-form-urlencoded; charset=utf-8",
	}

	if opt.trace != nil {
		opt.header[trace.Header] = []string{opt.trace.ID()}
	}

	ttl := opt.ttl
	if ttl <= 0 {
		ttl = DefaultTTL
	}

	// cancel ctx and release resource after ttl time
	ctx, cancel := context.WithTimeout(context.Background(), ttl)
	defer cancel()

	if opt.dialog != nil {
		decodedURL, _ := url.QueryUnescape(URL)
		opt.dialog.Request = &trace.Request{
			TTL:        ttl.String(),
			Method:     method,
			DecodedURL: decodedURL,
			Header:     opt.header,
		}
	}

	retryTimes := opt.retryTimes
	if retryTimes <= 0 {
		retryTimes = DefaultRetryTimes
	}

	retryDelay := opt.retryDelay
	if retryDelay <= 0 {
		retryDelay = DefaultRetryDelay
	}

	var httpCode int
	defer func() {
		// process alarm logic
		if opt.alarmObject == nil {
			return
		}

		if opt.alarmVerify != nil && !opt.alarmVerify(body) && err == nil {
			return
		}

		info := &struct {
			TraceId string `json:"trace_id"`
			Request struct {
				Method string `json:"method"`
				URL    string `json:"url"`
			} `json:"request"`
			Response struct {
				HttpCode int    `json:"http_code"`
				Body     string `json:"body"`
			} `json:"response"`
			Error string `json:"error"`
		}{}

		if opt.trace != nil {
			info.TraceId = opt.trace.ID()
		}

		info.Request.Method = method
		info.Request.URL = URL
		info.Response.HttpCode = httpCode
		info.Response.Body = string(body)
		info.Error = ""
		if err != nil {
			info.Error = fmt.Sprintf("%+v", err)
		}

		raw, _ := json.MarshalIndent(info, "", " ")
		onFailedAlarm(opt.alarmTitle, raw, opt.logger, opt.alarmObject)
	}()

	for k := 0; k < retryTimes; k++ {
		body, httpCode, err = doHTTP(ctx, method, URL, nil, opt)
		if shouldRetry(ctx, httpCode) || (opt.alarmVerify != nil && opt.retryVerify(body)) {
			time.Sleep(retryDelay)
			continue
		}
		return
	}
	return
}

// todo need to add withFormBody and withJsonBody logic

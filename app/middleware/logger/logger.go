package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-api-plus/app/config"
	"go-api-plus/app/utils/jsonutils"
	"go-api-plus/app/utils/responseutils"
	"go-api-plus/app/utils/timeutils"
	"log"
	"net/http"
	"os"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

var accessChannel = make(chan string, 100)

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// SetUp can set gin used middleware
func SetUp() gin.HandlerFunc {
	// let goroutine to process the log print
	go handleAccessChannel()

	return func(ctx *gin.Context) {
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		// that can let gin will do log store and then return responseutils
		ctx.Writer = bodyLogWriter

		// start time
		startTime := timeutils.GetCurrentMilliTime()

		// process request
		ctx.Next()

		responseBody := bodyLogWriter.body.String()

		var respCode int
		var respMsg string
		var respData interface{}

		if responseBody != "" {
			resp := responseutils.Response{}
			err := json.Unmarshal([]byte(responseBody), &resp)
			if err == nil {
				respCode = resp.Code
				respMsg = resp.Msg
				respData = resp.Data
			}
		}

		// end time
		entTime := timeutils.GetCurrentMilliTime()

		if ctx.Request.Method == http.MethodPost {
			_ = ctx.Request.ParseForm()
		}

		// log format
		accessFormat := make(map[string]interface{})

		accessFormat["request_time"] = startTime
		accessFormat["request_method"] = ctx.Request.Method
		accessFormat["request_uri"] = ctx.Request.RequestURI
		accessFormat["request_proto"] = ctx.Request.Proto
		accessFormat["request_user_agent"] = ctx.Request.UserAgent()
		accessFormat["request_post_data"] = ctx.Request.PostForm.Encode()
		accessFormat["request_client_ip"] = ctx.ClientIP()

		accessFormat["response_time"] = time.Now().Unix()
		accessFormat["response_code"] = respCode
		accessFormat["response_msg"] = respMsg
		accessFormat["response_data"] = respData

		accessFormat["elapse_time"] = fmt.Sprintf("%vms", entTime-startTime)

		accessLogJson, _ := jsonutils.JsonEncode(accessFormat)

		accessChannel <- accessLogJson
	}

}

func handleAccessChannel() {
	if f, err := os.OpenFile(config.AppAccessLogName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		log.Println(err)
	} else {
		for accessLog := range accessChannel {
			_, _ = f.WriteString(accessLog)
		}
	}
	return
}

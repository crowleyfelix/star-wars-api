package middlewares

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/aphistic/gomol"

	"github.com/gin-gonic/gin"
)

//Logger logs all requests and responses
func Logger(context *gin.Context) {
	rec := newInterceptor(context)
	var err string

	defer func() {
		gomol.Infom(gomol.NewAttrsFromMap(map[string]interface{}{
			"url":          context.Request.URL,
			"error":        err,
			"status":       rec.Status(),
			"durationTime": rec.Duration(),
			"request":      rec.Request(),
			"response":     rec.Response(),
		}), "request received")
	}()

	if e := rec.StartTrack(); e != nil {
		err = e.Error()
	}

	context.Next()
}

type recorder struct {
	gin.ResponseWriter
	startTime   time.Time
	requestBuf  *bytes.Buffer
	responseBuf *bytes.Buffer
	context     *gin.Context
}

func newInterceptor(context *gin.Context) *recorder {
	rec := &recorder{
		ResponseWriter: context.Writer,
		context:        context,
	}
	return rec
}

func (rec *recorder) Write(b []byte) (int, error) {
	if n, err := rec.responseBuf.Write(b); err != nil {
		return n, err
	}
	return rec.ResponseWriter.Write(b)
}

func (rec *recorder) StartTrack() error {
	rec.startTime = time.Now()
	rec.requestBuf = bytes.NewBufferString("")
	rec.responseBuf = bytes.NewBufferString("")

	rb, err := ioutil.ReadAll(rec.context.Request.Body)
	if err != nil {
		return err
	}

	reader := ioutil.NopCloser(bytes.NewBuffer(rb))
	rec.context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rb))

	if _, err = rec.requestBuf.ReadFrom(reader); err != nil {
		return err
	}

	rec.context.Writer = rec
	return nil
}

func (rec *recorder) Request() interface{} {
	var data interface{}
	if err := json.Unmarshal(rec.requestBuf.Bytes(), &data); err != nil {
		return rec.requestBuf.String()
	}

	return data
}

func (rec *recorder) Response() interface{} {
	var data interface{}
	if err := json.Unmarshal(rec.responseBuf.Bytes(), &data); err != nil {
		return nil
	}

	return data
}

func (rec *recorder) Duration() int64 {
	return time.Since(rec.startTime).Nanoseconds()
}

func (rec *recorder) Status() int {
	return rec.ResponseWriter.Status()
}

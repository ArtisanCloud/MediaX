package request

import (
	"bytes"
	"fmt"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"io"
	"net/http"
)

func LogRequest(needLog bool, logger *logger.Logger, request *http.Request) {
	if !needLog {
		return
	}
	var output bytes.Buffer
	// 前置中间件
	output.WriteString(fmt.Sprintf("%s %s ", request.Method, request.URL.String()))

	// print out request header
	output.Write([]byte("request header: { "))
	for k, vals := range request.Header {
		for _, v := range vals {
			output.Write([]byte(fmt.Sprintf("%s:%s", k, v)))
		}
	}
	output.Write([]byte("} "))

	// print out request body
	if request.Body != nil {

		output.Write([]byte("request body:"))
		var buf bytes.Buffer
		reader := io.TeeReader(request.Body, &buf)
		body, _ := io.ReadAll(reader)
		request.Body = io.NopCloser(&buf)
		output.Write(body)
	}

	logger.Info(output.String())
}

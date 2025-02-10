package response

import (
	"bytes"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"net/http"
	"net/http/httputil"
)

func LogResponse(needLog bool, l *logger.Logger, response *http.Response) {
	if !needLog {
		return
	}
	var output bytes.Buffer
	output.Write([]byte("------------------"))
	output.Write([]byte("response content:"))
	dumpRes, _ := httputil.DumpResponse(response, true)
	output.Write(dumpRes)
	l.Info(output.String())
}

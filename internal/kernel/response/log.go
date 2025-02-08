package response

import (
	"bytes"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"net/http"
	"net/http/httputil"
)

func LogResponse(l *logger.Logger, response *http.Response) {
	var output bytes.Buffer
	output.Write([]byte("------------------"))
	output.Write([]byte("response content:"))
	dumpRes, _ := httputil.DumpResponse(response, true)
	output.Write(dumpRes)
	l.Info(output.String())
}

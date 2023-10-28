package helpers

import (
	"io"
	"net/http"
)

func calculateDataSent(req *http.Request) int64 {
	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err == nil {
			return int64(len(body))
		}
	}
	return 0
}

func calculateDataReceived(resp *http.Response) int64 {
	if resp.Body != nil {
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			return int64(len(body))
		}
	}
	return 0
}

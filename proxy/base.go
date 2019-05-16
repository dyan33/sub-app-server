package proxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func NewResponse(r *http.Request, contentType string, status int, body []byte) *http.Response {
	resp := &http.Response{}
	resp.Request = r
	resp.TransferEncoding = r.TransferEncoding
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Type", contentType)
	resp.StatusCode = status

	buf := bytes.NewBuffer(body)

	resp.ContentLength = int64(buf.Len())
	resp.Body = ioutil.NopCloser(buf)
	return resp
}

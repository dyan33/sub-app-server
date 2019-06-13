package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type HttpRequest struct {
	Id     int64             `json:"id"`
	Method string            `json:"method"`
	Url    string            `json:"url"`
	Header map[string]string `json:"header"`
	Body   []byte            `json:"body"`
}

type HttpResponse struct {
	Id      int64       `json:"id"`
	Code    int         `json:"code"`
	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`
}

func makeRequest(id int64, req *http.Request) *HttpRequest {

	body, _ := ioutil.ReadAll(req.Body)

	request := &HttpRequest{
		Id:     id,
		Method: req.Method,
		Header: map[string]string{},
		Url:    req.URL.String(),
		Body:   body,
	}

	for key, value := range req.Header {

		if key == "Proxy-Connection" {
			continue
		}

		request.Header[key] = value[0]
	}

	return request
}

func makeResponse(r *http.Request, response HttpResponse) *http.Response {

	buf := bytes.NewBuffer(response.Body)

	return &http.Response{
		Request:          r,
		TransferEncoding: r.TransferEncoding,
		Header:           response.Headers,
		StatusCode:       response.Code,
		ContentLength:    int64(buf.Len()),
		Body:             ioutil.NopCloser(buf),
	}

}

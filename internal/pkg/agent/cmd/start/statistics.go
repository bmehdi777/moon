package start

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type HttpCall struct {
	Request  http.Request
	Response http.Response
}
type HttpCalls []HttpCall

type Statistics struct {
	HttpCalls HttpCalls
	Event     chan int
}

type RequestMessage struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    string `json:"body"`
}

type ResponseMessage struct {
	Status  int `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string `json:"body"`
}

type HttpMessage struct {
	Request  RequestMessage `json:"request"`
	Response ResponseMessage `json:"response"`
}

func (call *HttpCall) ToHttpMessage() HttpMessage {
	defer call.Request.Body.Close()
	defer call.Response.Body.Close()

	var reqBody string
	reqBodyBytes, err := io.ReadAll(call.Request.Body)
	if err != nil {
		reqBody = ""
	}
	reqBody = string(reqBodyBytes)

	var respBody string
	respBodyBytes, err := io.ReadAll(call.Response.Body)
	if err != nil {
		respBody = ""
	}
	respBody = string(respBodyBytes)


	reqHeaders := make(map[string]string, 0)
	for name, values := range call.Request.Header {
		reqHeaders[name] = strings.Join(values, ", ")
	}

	respHeaders := make(map[string]string, 0)
	for name, values := range call.Response.Header {
		respHeaders[name] = strings.Join(values, ", ")
	}

	return HttpMessage{
		Request: RequestMessage{
			Method:  call.Request.Method,
			Path:    call.Request.URL.Path,
			Headers: reqHeaders,
			Body:    reqBody,
		},
		Response: ResponseMessage{
			Status:  call.Response.StatusCode,
			Headers: respHeaders,
			Body:    respBody,
		},
	}
}


func (h HttpCalls) Format() string {
	msgs := make([]HttpMessage, 0)

	for _, elem := range h {
		msgs = append(msgs, elem.ToHttpMessage())
	}

	msgsBytes, err := json.Marshal(msgs)
	if err != nil {
		return ""
	}

	return string(msgsBytes)
}

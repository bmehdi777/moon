package start

import (
	"encoding/json"
)

type HttpMessages []HttpMessage

type Statistics struct {
	HttpCalls     HttpMessages
	Event         chan int
	EventListener int
}

type RequestMessage struct {
	Method   string            `json:"method"`
	Path     string            `json:"path"`
	Headers  map[string]string `json:"headers"`
	Body     string            `json:"body"`
	Datetime string            `json:"datetime"`
}

type ResponseMessage struct {
	Status   int               `json:"status"`
	Headers  map[string]string `json:"headers"`
	Body     string            `json:"body"`
	Duration string            `json:"duration"`
}

type HttpMessage struct {
	Request  RequestMessage  `json:"request"`
	Response ResponseMessage `json:"response"`
}

func (h HttpMessages) Format() string {
	msgs := make([]HttpMessage, 0)

	for _, elem := range h {
		msgs = append(msgs, elem)
	}

	msgsBytes, err := json.Marshal(msgs)
	if err != nil {
		return ""
	}

	return string(msgsBytes)
}

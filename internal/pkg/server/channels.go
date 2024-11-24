package server

import "net/http"

type ChannelsHttp struct {
	RequestChannel  chan *http.Request
	ResponseChannel chan *http.Response
}

type ChannelsDomains map[string]ChannelsHttp

func (c *ChannelsDomains) Add(name string) {
	(*c)[name] = ChannelsHttp{
		RequestChannel:  make(chan *http.Request),
		ResponseChannel: make(chan *http.Response),
	}
}

func (c *ChannelsDomains) Delete(name string) {
	delete(*c, name)
}

func (c *ChannelsDomains) Get(name string) *ChannelsHttp {
	if channel, found := (*c)[name]; found {
		return &channel
	}
	return nil
}

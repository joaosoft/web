package client

import "webserver"

func (client *WebClient) NewRequest() (*Request, error) {

	request := &Request{
		Base: Base{
			Headers:   make(webserver.Headers),
			Cookies:   make(webserver.Cookies),
			Params:    make(webserver.Params),
			UrlParams: make(webserver.UrlParams),
			Charset:   webserver.CharsetUTF8,
			client:    client,
		},
		Attachments: make(map[string]Attachment),
	}

	return request, nil
}

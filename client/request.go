package client

import "web"

func (client *WebClient) NewRequest() (*Request, error) {

	request := &Request{
		Base: Base{
			Headers:   make(web.Headers),
			Cookies:   make(web.Cookies),
			Params:    make(web.Params),
			UrlParams: make(web.UrlParams),
			Charset:   web.CharsetUTF8,
			client:    client,
		},
		Attachments: make(map[string]Attachment),
	}

	return request, nil
}

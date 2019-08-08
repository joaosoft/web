package web

import (
	"fmt"
	"regexp"
)

func (c *Client) NewRequest(method Method, url string) (*Request, error) {

	// validate url
	regx := regexp.MustCompile(RegexForURL)
	if !regx.MatchString(url) {
		return nil, fmt.Errorf("invalid url [%s]", url)
	}

	address := NewAddress(url)
	params := address.Params

	return &Request{
		Base: Base{
			Client:      c,
			Protocol:    ProtocolHttp1p1,
			Method:      method,
			Address:     address,
			Headers:     make(Headers),
			Cookies:     make(Cookies),
			Params:      params,
			UrlParams:   make(UrlParams),
			Charset:     CharsetUTF8,
			ContentType: ContentTypeApplicationJSON,
		},
		FormData:            make(map[string]*FormData),
		Attachments:         make(map[string]*Attachment),
		MultiAttachmentMode: c.multiAttachmentMode,
		Boundary:            RandomBoundary(),
	}, nil
}

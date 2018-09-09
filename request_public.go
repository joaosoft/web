package web

func (r *Request) SetHeader(name string, header Header) {
	r.Headers[HeaderType(name)] = header
}

func (r *Request) GetHeader(name string) *Header {
	if header, ok := r.Headers[HeaderType(name)]; ok {
		return &header
	}

	return nil
}

func (r *Request) SetCookie(name string, cookie Cookie) {
	r.Cookies[name] = cookie
}

func (r *Request) GetCookie(name string) *Cookie {
	if cookie, ok := r.Cookies[name]; ok {
		return &cookie
	}

	return nil
}

func (r *Request) SetParam(name string, param Param) {
	r.Params[name] = param
}

func (r *Request) GetParam(name string) *Param {
	if param, ok := r.Params[name]; ok {
		return &param
	}

	return nil
}

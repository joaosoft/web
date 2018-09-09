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

func (r *Request) SetParm(name string, parm Parm) {
	r.Parms[name] = parm
}

func (r *Request) GetParm(name string) *Parm {
	if parm, ok := r.Parms[name]; ok {
		return &parm
	}

	return nil
}

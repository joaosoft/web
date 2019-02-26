package web

import "strings"

type Params map[string][]string

type Address struct {
	Full   string
	Url    string
	Host   string
	Params Params
}

func NewAddress(url string) *Address {
	address := &Address{
		Full:   url,
		Url:    url,
		Host:   strings.SplitN(url, "/", 2)[0],
		Params: make(Params),
	}

	// load query parameters
	if split := strings.SplitN(url, "?", 2); len(split) > 1 {
		address.Url = string(split[0])
		if parms := strings.Split(split[1], "&"); len(parms) > 0 {
			for _, parm := range parms {
				if p := strings.Split(parm, "="); len(p) > 1 {
					if split := strings.SplitN(p[1], ",", -1); len(split) > 0 {
						for _, s := range split {
							address.Params[p[0]] = append(address.Params[p[0]], s)
						}
					}
					address.Params[p[0]] = append(address.Params[p[0]], p[1])
				}
			}
		}
	}

	return address
}

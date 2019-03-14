package web

import "strings"

type Params map[string][]string

type Address struct {
	Full   string
	Schema string
	Url    string
	Host   string
	Params Params
}

func NewAddress(url string) *Address {
	var tmp, full, schema, host string
	var params = make(Params)

	tmp = url
	full = tmp // full

	split := strings.SplitN(tmp, "//", 2)
	if len(split) == 2 {
		schema = split[0] // schema
		tmp = split[1]
	}

	split = strings.SplitN(tmp, "/", 2)
	host = split[0] // host

	if len(split) == 2 {
		tmp = split[1]
		url = tmp // url
	}

	// load query parameters
	if split := strings.SplitN(tmp, "?", 2); len(split) > 1 {
		url = string(split[0]) // url
		if parms := strings.Split(split[1], "&"); len(parms) > 0 {
			for _, parm := range parms {
				if p := strings.Split(parm, "="); len(p) > 1 {
					if split := strings.SplitN(p[1], ",", -1); len(split) > 0 {
						for _, s := range split {
							params[p[0]] = append(params[p[0]], s)
						}
					}
					params[p[0]] = append(params[p[0]], p[1])
				}
			}
		}
	}

	return &Address{
		Full:   full,
		Schema: schema,
		Host:   host,
		Url:    url,
		Params: params,
	}
}

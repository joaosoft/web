package web

import "fmt"

func (ctx *Context) Redirect(host string) error {
	if ctx.Request.Client == nil {
		client, err := NewClient(WithClientLogger(ctx.Request.Server.logger))
		if err != nil {
			return err
		}

		ctx.Request.Client = client
	}

	ctx.Request.Address = NewAddress(fmt.Sprintf("%s%s", host, ctx.Request.Address.Url))

	response, err := ctx.Request.Send()
	if err != nil {
		return err
	}

	return ctx.Response.Bytes(response.Status, response.ContentType, response.Body)
}

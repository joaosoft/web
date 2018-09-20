package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"web"

	"github.com/labstack/gommon/log"
)

func (w *WebServer) handlerFile(ctx *Context) error {
	log.Infof("handling file %s", ctx.Request.FullUrl)

	dir, _ := os.Getwd()
	path := fmt.Sprintf("%s%s", dir, ctx.Request.FullUrl)

	if _, err := os.Stat(path); err == nil {
		if bytes, err := ioutil.ReadFile(path); err != nil {
			ctx.Response.Status = web.StatusNotFound
		} else {
			ctx.Response.Status = web.StatusOK
			ctx.Response.Body = bytes
		}
	} else {
		ctx.Response.Status = web.StatusNotFound
	}

	return nil
}

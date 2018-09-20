package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"web/common"

	"github.com/labstack/gommon/log"
)

func (w *Server) handlerFile(ctx *Context) error {
	log.Infof("handling file %s", ctx.Request.FullUrl)

	dir, _ := os.Getwd()
	path := fmt.Sprintf("%s%s", dir, ctx.Request.FullUrl)

	if _, err := os.Stat(path); err == nil {
		if bytes, err := ioutil.ReadFile(path); err != nil {
			ctx.Response.Status = common.StatusNotFound
		} else {
			ctx.Response.Status = common.StatusOK
			ctx.Response.Body = bytes
		}
	} else {
		ctx.Response.Status = common.StatusNotFound
	}

	return nil
}

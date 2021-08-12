package web

import (
	"net/http"

	"errors"
)

var (
	ErrorInvalidAuthorization = errors.New(errors.LevelError, http.StatusUnauthorized, "invalid authorization")
)

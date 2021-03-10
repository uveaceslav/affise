package controller

import (
	"net/http"
)

type MainController interface {
	Index(rw http.ResponseWriter, r *http.Request)
}

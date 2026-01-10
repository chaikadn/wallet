package http

import (
	"net/http"
	"wallet/internal/service"
)

type HttpHandler struct {
	sr service.Service
}

func HandlePostWallet(w http.ResponseWriter, r *http.Request) {

}

func HandleGetWallet(w http.ResponseWriter, r *http.Request) {

}

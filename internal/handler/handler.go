package handler

import "net/http"

type Handler interface {
	HandlePostWallet(w http.ResponseWriter, r *http.Request)
	HandleGetWallet(w http.ResponseWriter, r *http.Request)
}

package web

import "github.com/go-chi/chi"

func NewRouter(svr *Server) *chi.Mux {
	r := chi.NewRouter()
	r.Use(svr.Auth)
	r.Post("/oss", svr.Upload)
	r.Post("/oss/{provider}", svr.Upload)
	r.Get("/oss", svr.Handle)
	r.Get("/oss/{provider}", svr.Handle)
	return r
}

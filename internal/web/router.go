package web

import "github.com/go-chi/chi"

func NewRouter(svr *Server) *chi.Mux {
	r := chi.NewRouter()
	r.Use(svr.Auth)
	r.Get("/oss/{provider}", svr.Handle)
	return r
}

package server

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pzl/mstk"
	"github.com/pzl/mstk/logger"
)

type StaticHandler http.Handler

func (s *Server) SetupRoutes() {
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RequestLogger(logger.NewChi(s.Log)))
	s.Router.Use(middleware.Heartbeat("/ping"))
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(cors)

	s.routeAPI()
	s.routeWeb()
}

func (s *Server) routeWeb() {
	s.Router.Get("/_nuxt/*", s.AssetHandler.ServeHTTP)
	// add other "/" root level static files needed here
	files := []string{"favicon.ico"}
	for _, f := range files {
		s.Router.Get("/"+f, s.AssetHandler.ServeHTTP)
	}

	// the above probably not necessary if we route the rest to vue
	s.Router.Get("/*", s.AssetHandler.ServeHTTP)
}

func (s *Server) routeAPI() {
	s.Router.Route("/api/v1", func(v1 chi.Router) {
		v1.Use(contentJSON) // by default
		v1.Use(mstk.APIVer(1))

		v1.Get("/file", s.FileListHandler)
		v1.Post("/file", s.FileCreateHandler)
		v1.Get("/file/{id}", s.FileGetHandler)
		v1.Put("/file/{id}", s.FileCreateManualHandler)
		v1.Delete("/file/{id}", s.FileDeleteHandler)

		v1.Get("/link", s.LinkListHandler)
		v1.Post("/link", s.LinkCreateHandler)
		v1.Get("/link/{id}", s.LinkGetHandler)
		v1.Put("/link/{id}", s.LinkCreateManualHandler)
		v1.Delete("/link/{id}", s.LinkDeleteHandler)

		v1.Get("/text", s.TextListHandler)
		v1.Post("/text", s.TextCreateHandler)
		v1.Get("/text/{id}", s.TextGetHandler)
		v1.Put("/text/{id}", s.TextCreateManualHandler)
		v1.Delete("/text/{id}", s.TextDeleteHandler)
	})
}

func contentJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			// OPTIONS -- let's handle fully in here

			w.Header().Add("Vary", "Origin")
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Add("Vary", "Access-Control-Request-Headers")

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", strings.ToUpper(r.Header.Get("Access-Control-Request-Method")))
			w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, User, Content-Length, Accept-Encoding, X-CSRF-Token")

			w.WriteHeader(http.StatusOK)
			return
			// and be done with it. DO NOT serveHTTP
		}
		// not OPTIONS, do some header alteration and pass on
		w.Header().Add("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, User, Content-Length, Accept-Encoding, X-CSRF-Token")
		next.ServeHTTP(w, r)

	})
}

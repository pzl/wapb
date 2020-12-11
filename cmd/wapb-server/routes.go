package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pzl/mstk"
	"github.com/pzl/mstk/logger"
	"github.com/sirupsen/logrus"
)

func makeRouter(log *logrus.Logger, cfg Config) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Use(middleware.RequestLogger(logger.NewChi(log)))
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.Recoverer)
	router.Use(cors)

	// setup server static asset handling
	handleWeb(router, cfg, log)

	handleAPI(router, cfg, log)

	return router
}

func handleWeb(router *chi.Mux, cfg Config, log *logrus.Logger) {
	var assetHandler http.Handler
	if cfg.DevMode {
		log.Info("dev mode enabled. Listening to npm dev server at localhost:3000")
		devServer, _ := url.Parse("http://localhost:3000")
		proxy := httputil.NewSingleHostReverseProxy(devServer)
		assetHandler = proxy
	} else {
		log.Info("serving in production mode, with precompiled assets")
		assetHandler = http.FileServer(assets)
	}
	router.Get("/_nuxt/*", func(w http.ResponseWriter, r *http.Request) {
		assetHandler.ServeHTTP(w, r)
	})
	// add other "/" root level static files needed here
	files := []string{"favicon.ico"}
	for _, f := range files {
		router.Get("/"+f, func(w http.ResponseWriter, r *http.Request) {
			assetHandler.ServeHTTP(w, r)
		})
	}
}

func handleAPI(router *chi.Mux, cfg Config, log *logrus.Logger) {

	router.Route("/api/v1", func(v1 chi.Router) {
		v1.Use(contentJSON)
		v1.Use(mstk.APIVer(1))

		v1.Get("/file", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{ "data": [] }`))
		})
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
			w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, User")

			w.WriteHeader(http.StatusOK)
			// and be done with it. DO NOT serveHTTP
		} else {
			// not OPTIONS, do some header alteration and pass on
			w.Header().Add("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, User")
			next.ServeHTTP(w, r)
		}
	})
}

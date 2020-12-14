package server

import (
	"net/http"
)

type File struct {
	ID string
}

func (s *Server) FileListHandler(w http.ResponseWriter, r *http.Request) {}
func (s *Server) FileGetHandler(w http.ResponseWriter, r *http.Request)  {}
func (s *Server) FileCreateHandler(w http.ResponseWriter, r *http.Request) {
	f := File{}
	if err := jsCfg.NewDecoder(r.Body).Decode(&f); err != nil {
		s.Log.WithError(err).Error("error json decoding body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
func (s *Server) FileCreateManualHandler(w http.ResponseWriter, r *http.Request) {}
func (s *Server) FileDeleteHandler(w http.ResponseWriter, r *http.Request)       {}

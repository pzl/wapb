package server

import "net/http"

func (s *Server) LinkListHandler(w http.ResponseWriter, r *http.Request)         {}
func (s *Server) LinkGetHandler(w http.ResponseWriter, r *http.Request)          {}
func (s *Server) LinkCreateHandler(w http.ResponseWriter, r *http.Request)       {}
func (s *Server) LinkCreateManualHandler(w http.ResponseWriter, r *http.Request) {}
func (s *Server) LinkDeleteHandler(w http.ResponseWriter, r *http.Request)       {}

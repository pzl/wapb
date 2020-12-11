package server

import (
	"bytes"
	"net/http"
)

type File struct {
	ID string
}

func (s *Server) FileListHandler(w http.ResponseWriter, r *http.Request) {
	files := make([][]byte, 0, 20)

	err := getAllForType(s.DB, StorageFileKey, func(v []byte) error {
		buf := make([]byte, len(v))
		copy(buf, v)
		files = append(files, buf)
		return nil
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	pfx := []byte(`{"data":[`)
	out := bytes.Join(files, []byte{','})

	out = append(pfx, out...)
	out = append(out, []byte(`]}`)...)
	w.Write(out)

}
func (s *Server) FileGetHandler(w http.ResponseWriter, r *http.Request) {}
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

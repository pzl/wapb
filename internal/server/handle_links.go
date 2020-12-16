package server

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dgraph-io/badger/v2"
	"github.com/go-chi/chi"
)

type Link struct {
	CommonFields
	URL string `json:"url"`
}

func (s *Server) LinkListHandler(w http.ResponseWriter, r *http.Request) {
	s.doListHandler(w, StorageLinkKey)
}
func (s *Server) LinkGetHandler(w http.ResponseWriter, r *http.Request) {
	handlers := map[string]func([]byte){
		"text/plain": func(buf []byte) {
			var l Link
			if err := jsCfg.Unmarshal(buf, &l); err != nil {
				s.Log.WithError(err).Error("error unserializing record")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(l.URL + "\n"))
		},
	}

	s.doGetOneHandler(w, r, StorageLinkKey, handlers)
}
func (s *Server) LinkCreateHandler(w http.ResponseWriter, r *http.Request) {
	var cr Link

	handlers := map[string]CreateHandlerFunc{
		"text/plain": func(body io.Reader, v url.Values) error {
			buf, err := ioutil.ReadAll(body)
			cr.URL = string(buf)
			return err
		},
		"application/x-www-form-urlencoded": func(body io.Reader, v url.Values) error {
			cr.URL = v.Get("url")
			return nil
		},
		"application/json": func(body io.Reader, v url.Values) error {
			return jsCfg.NewDecoder(body).Decode(&cr)
		},
	}

	ct, err := s.doCreateHandler(r, &cr.CommonFields, handlers)
	if err != nil {
		s.Log.WithError(err).Error("error creating record")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if cr.URL == "" {
		s.Log.Debug("ignoring empty string upload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	buf, err := jsCfg.Marshal(cr)
	if err != nil {
		s.Log.WithError(err).Error("error serializing link record")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := writeBytes(s.DB, StorageLinkKey, cr.ID, buf, makeMeta(cr.CommonFields), cr.TTL); err != nil {
		s.Log.WithError(err).Error("error writing link record")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if accept not specific, then send back the same format we got
	accept := r.Header.Get("Accept")
	if accept == "" || accept == "*/*" {
		accept = ct
	}

	if accept == "text/plain" {
		w.Header().Set("Content-Type", accept) // any header changes must happen BEFORE WriteHeader
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://" + r.Host + "/link/" + cr.ID + "\n"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(buf)
}

func (s *Server) LinkCreateManualHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (s *Server) LinkDeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := deleteRecord(s.DB, StorageLinkKey, chi.URLParam(r, "id"))
	if err == badger.ErrKeyNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

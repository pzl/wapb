package server

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dgraph-io/badger/v2"
	"github.com/go-chi/chi"
)

type Text struct {
	CommonFields
	Text string `json:"text"`
}

func (s *Server) TextListHandler(w http.ResponseWriter, r *http.Request) {
	s.doListHandler(w, StorageTextKey)
}
func (s *Server) TextGetHandler(w http.ResponseWriter, r *http.Request) {
	handlers := map[string]func([]byte){
		"text/plain": func(buf []byte) {
			var t Text
			if err := jsCfg.Unmarshal(buf, &t); err != nil {
				s.Log.WithError(err).Error("error unserializing record")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(t.Text))
		},
	}

	s.doGetOneHandler(w, r, StorageTextKey, handlers)
}
func (s *Server) TextCreateHandler(w http.ResponseWriter, r *http.Request) {
	var cr Text

	handlers := map[string]CreateHandlerFunc{
		"text/plain": func(body io.Reader, v url.Values) error {
			buf, err := ioutil.ReadAll(body)
			cr.Text = string(buf)
			return err
		},
		"application/x-www-form-urlencoded": func(body io.Reader, v url.Values) error {
			cr.Text = v.Get("text")
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

	if cr.Text == "" {
		s.Log.Debug("ignoring empty string upload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	buf, err := jsCfg.Marshal(cr)
	if err != nil {
		s.Log.WithError(err).Error("error serializing text record")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := writeBytes(s.DB, StorageTextKey, cr.ID, buf, makeMeta(cr.CommonFields), cr.TTL); err != nil {
		s.Log.WithError(err).Error("error writing text record")
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
		w.Write([]byte("http://" + r.Host + "/text/" + cr.ID + "\n"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(buf)
}
func (s *Server) TextCreateManualHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (s *Server) TextDeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := deleteRecord(s.DB, StorageTextKey, chi.URLParam(r, "id"))
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

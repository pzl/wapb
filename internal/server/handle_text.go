package server

import (
	"io/ioutil"
	"net/http"

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

	ct, body := getContentType(r)

	if ct == "text/plain" {
		buf, err := ioutil.ReadAll(body)
		if err != nil {
			s.Log.WithError(err).Error("error reading text body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		setCommonFieldsByQuery(&cr.CommonFields, r)
		cr.Text = string(buf)
	} else {
		if err := jsCfg.NewDecoder(body).Decode(&cr); err != nil {
			s.Log.WithError(err).Error("error decoding text create body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if cr.Text == "" {
		s.Log.Debug("ignoring empty string upload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// overwrite non-user-providable fields
	setCreateCommonFields(&cr.CommonFields)

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

	// if accept not given, then send back the same format we got
	accept := r.Header.Get("Accept")
	if accept == "" {
		accept = ct
	}

	if accept == "text/plain" {
		w.Header().Set("Content-Type", accept) // any header changes must happen BEFORE WriteHeader
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://" + r.Host + "/text/" + cr.ID))
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

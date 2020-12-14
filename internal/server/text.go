package server

import (
	"bytes"
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

	files, err := getAllForType(s.DB, StorageTextKey, func(v []byte) ([]byte, error) {
		var t Text
		err := jsCfg.Unmarshal(v, &t)
		if err != nil {
			return nil, err
		}
		t.Text = ""
		return jsCfg.Marshal(t)

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
func (s *Server) TextGetHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	buf, err := getOneBytes(s.DB, StorageTextKey, id)
	if err != nil && err == badger.ErrKeyNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		s.Log.WithError(err).Error("error fetching record")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Header.Get("Accept") == "text/plain" {
		var t Text
		if err := jsCfg.Unmarshal(buf, &t); err != nil {
			s.Log.WithError(err).Error("error unserializing record")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(t.Text))
		return
	}

	w.Write(buf)

}
func (s *Server) TextCreateHandler(w http.ResponseWriter, r *http.Request) {
	var cr struct {
		CommonFields
		Text string `json:"text"`
	}

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
	t := Text{
		cr.CommonFields,
		cr.Text,
	}

	buf, err := jsCfg.Marshal(t)
	if err != nil {
		s.Log.WithError(err).Error("error serializing text record")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := writeBytes(s.DB, StorageTextKey, t.ID, buf, makeMeta(t.CommonFields), t.TTL); err != nil {
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
		w.Write([]byte("http://" + r.Host + "/text/" + t.ID))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(buf)

}
func (s *Server) TextCreateManualHandler(w http.ResponseWriter, r *http.Request) {}
func (s *Server) TextDeleteHandler(w http.ResponseWriter, r *http.Request)       {}

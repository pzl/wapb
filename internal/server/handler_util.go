package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/go-chi/chi"
)

type CommonFields struct {
	BurnAfterRead bool   `json:"burn,omitempty"`
	Hidden        bool   `json:"hidden,omitempty"`
	TTL           int64  `json:"ttl,omitempty"`
	ID            string `json:"id,omitempty"`
	Created       int64  `json:"created,omitempty"` // timestamp of creation

}

func (s *Server) doListHandler(w http.ResponseWriter, sk StorageKey) {
	items, err := getAllForType(s.DB, sk)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	pfx := []byte(`{"data":[`)
	out := bytes.Join(items, []byte{','})

	out = append(pfx, out...)
	out = append(out, []byte(`]}`)...)
	w.Write(out)
}

func (s *Server) doGetOneHandler(w http.ResponseWriter, r *http.Request, sk StorageKey, handlers map[string]func([]byte)) {
	id := chi.URLParam(r, "id")

	buf, err := getOneBytes(s.DB, nil, sk, id)
	if err == badger.ErrKeyNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		s.Log.WithError(err).Error("error fetching record")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// if a custom handler was passed, respond with that Otherwise parrot out the bytes
	if handler, exists := handlers[r.Header.Get("Accept")]; exists {
		handler(buf)
		return
	}

	w.Write(buf)
}

type CreateHandlerFunc func(io.Reader, url.Values) error

func (s *Server) doCreateHandler(r *http.Request, c *CommonFields, handlers map[string]CreateHandlerFunc) (string, error) {
	values := r.URL.Query()
	ct, body := getContentType(r)
	s.Log.WithField("content-type", ct).Debug("detected content type")

	// pre-process some common fields if sent via form values
	if ct == "application/x-www-form-urlencoded" {
		b, err := ioutil.ReadAll(body)
		if err != nil {
			return ct, err
		}
		body = bytes.NewReader(b) // duplicate what we read, to send to the handler

		// get form data
		values, err = url.ParseQuery(string(b))
		if err != nil {
			return ct, err
		}
		qvals, err := url.ParseQuery(r.URL.RawQuery) // get query params
		if err != nil {
			return ct, err
		}

		copyValues(values, qvals) // merge query over form data
	}
	setCommonFieldsByValues(c, values)

	if handler, exists := handlers[ct]; exists {
		if err := handler(body, values); err != nil {
			return ct, err
		}
	}

	// overwrite non-user-providable fields
	setCreateCommonFields(c)

	return ct, nil
}

func censorPreventBurn(sk StorageKey, data []byte) ([]byte, error) {
	switch sk {
	case StorageFileGroupKey:
		fg := FileGroup{}
		if err := jsCfg.Unmarshal(data, &fg); err != nil {
			return nil, err
		}
		fg.Files = nil
		return jsCfg.Marshal(fg)
	case StorageTextKey:
		t := Text{}
		if err := jsCfg.Unmarshal(data, &t); err != nil {
			return nil, err
		}
		t.Text = ""
		return jsCfg.Marshal(t)
	case StorageLinkKey:
		l := Link{}
		if err := jsCfg.Unmarshal(data, &l); err != nil {
			return nil, err
		}
		return jsCfg.Marshal(l)
	case StorageFileKey:
		return nil, errors.New("unable to uncensor file contents")
	}

	return nil, errors.New("type not found")
}

func setCommonFieldsByValues(c *CommonFields, r url.Values) {
	switch strings.ToLower(r.Get("burn")) {
	case "true", "1", "yes", "y", "t":
		c.BurnAfterRead = true
	default:
		c.BurnAfterRead = false
	}

	switch ttl, err := strconv.ParseInt(r.Get("ttl"), 10, 64); {
	case err != nil:
		c.TTL = 0
	case ttl <= 0:
		c.TTL = 0
	default:
		c.TTL = ttl
	}

	switch strings.ToLower(r.Get("hidden")) {
	case "true", "1", "yes", "y", "t":
		c.Hidden = true
	default:
		c.Hidden = false
	}
}

func setCreateCommonFields(c *CommonFields) {
	c.ID = newID()
	c.Created = time.Now().Unix()
}

func getContentType(r *http.Request) (string, io.Reader) {
	bufd := bufio.NewReader(r.Body)

	ct := r.Header.Get("Content-Type")
	ct, _, _ = mime.ParseMediaType(ct) // strip off optional params & extras

	// this means values should be of form
	// key=value&key2=value2
	// but `curl` sends this header by default for ANY -d usage
	if ct == "application/x-www-form-urlencoded" {
		pk, err := bufd.Peek(15)
		if err != nil && err != io.EOF {
			return ct, bufd
		}
		if bytes.ContainsRune(pk, '=') {
			// it probably is k=v values
			return ct, bufd
		}

		// at this point, it probably is NOT form encoded data
		// so let's give json or text detection a try below
		ct = ""
	}

	if ct == "" {
		pk, err := bufd.Peek(5)
		if err != nil && err != io.EOF {
			// header not given, couldn't peek for some reason, assume json
			return "application/json", bufd
		}
		trim := bytes.TrimSpace(pk)
		if len(trim) > 0 && trim[0] == '{' {
			ct = "application/json"
		} else {
			ct = "text/plain"
		}

	}
	return ct, bufd
}

func newID() string {
	rand.Seed(time.Now().UnixNano())
	r := make([]byte, 4)
	n, err := rand.Read(r)
	if err != nil || n < 4 {
		// odd. OK let's do it an error-free way
		binary.LittleEndian.PutUint64(r, uint64(time.Now().Unix()))
	}
	return hex.EncodeToString(r[:3])
}

// https://golang.org/src/net/http/request.go#L1166
func copyValues(dst, src url.Values) {
	for k, vs := range src {
		dst[k] = append(dst[k], vs...)
	}
}

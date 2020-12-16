package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"math/rand"
	"net/http"
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

func (s *Server) doGetOneHandler(w http.ResponseWriter, r *http.Request, sk StorageKey, doers map[string]func([]byte)) {
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
	if cb, has := doers[r.Header.Get("Accept")]; has {
		cb(buf)
		return
	}

	w.Write(buf)
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

func setCommonFieldsByQuery(c *CommonFields, r *http.Request) {
	switch strings.ToLower(r.URL.Query().Get("burn")) {
	case "true", "1", "yes", "y", "t":
		c.BurnAfterRead = true
	default:
		c.BurnAfterRead = false
	}

	switch ttl, err := strconv.ParseInt(r.URL.Query().Get("ttl"), 10, 64); {
	case err != nil:
		c.TTL = 0
	case ttl <= 0:
		c.TTL = 0
	default:
		c.TTL = ttl
	}

	switch strings.ToLower(r.URL.Query().Get("hidden")) {
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
	if ct == "" {
		pk, err := bufd.Peek(5)
		if err != nil {
			// header not given, couldn't peek for some reason, assume json
			ct = "application/json"
		} else {
			trim := bytes.TrimSpace(pk)
			if len(trim) > 0 && trim[0] == '{' {
				ct = "application/json"
			} else {
				ct = "text/plain"
			}
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

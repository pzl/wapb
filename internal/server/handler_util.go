package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CommonFields struct {
	BurnAfterRead bool   `json:"burn,omitempty"`
	Hidden        bool   `json:"hidden,omitempty"`
	TTL           int    `json:"ttl,omitempty"`
	ID            string `json:"id,omitempty"`
	Created       int64  `json:"created,omitempty"` // timestamp of creation

}

func setCommonFieldsByQuery(c *CommonFields, r *http.Request) {
	switch strings.ToLower(r.URL.Query().Get("burn")) {
	case "true", "1", "yes", "y", "t":
		c.BurnAfterRead = true
	default:
		c.BurnAfterRead = false
	}

	switch ttl, err := strconv.Atoi(r.URL.Query().Get("ttl")); {
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
	return hex.EncodeToString(r[:2])
}

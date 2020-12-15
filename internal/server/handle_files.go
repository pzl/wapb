package server

import (
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/dgraph-io/badger/v2"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type FileGroup struct {
	CommonFields
	Files []File `json:"files,omitempty"`
}

type File struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	FileName string `json:"filename"`
	Mime     string `json:"mime"`
}

func (s *Server) FileGroupListHandler(w http.ResponseWriter, r *http.Request) {
	s.doListHandler(w, StorageFileGroupKey)
}
func (s *Server) FileGroupGetHandler(w http.ResponseWriter, r *http.Request) {
	s.doGetOneHandler(w, r, StorageFileGroupKey, nil)
}

func (s *Server) FileGroupCreateHandler(w http.ResponseWriter, r *http.Request) {
	var cr FileGroup

	ct, body := getContentType(r)
	if err := jsCfg.NewDecoder(body).Decode(&cr); err != nil && err != io.EOF {
		s.Log.WithError(err).Error("error decoding filegroup create body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// overwrite non-user-providable fields
	setCreateCommonFields(&cr.CommonFields)
	cr.Files = nil

	buf, err := jsCfg.Marshal(cr)
	if err != nil {
		s.Log.WithError(err).Error("error serializing filegroup record")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := writeBytes(s.DB, StorageFileGroupKey, cr.ID, buf, makeMeta(cr.CommonFields), cr.TTL); err != nil {
		s.Log.WithError(err).Error("error writing filegroup record")
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
		w.Write([]byte("http://" + r.Host + "/file/" + cr.ID))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(buf)
}
func (s *Server) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var fg FileGroup
	if err := getOne(s.DB, DontBurn, StorageFileGroupKey, id, &fg); err != nil {
		if err == badger.ErrKeyNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.Log.WithError(err).Error("error getting filegroup record")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	meta, err := getMeta(s.DB, StorageFileGroupKey, id)
	if err != nil {
		s.Log.WithError(err).Error("error getting filegroup meta")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reader, err := r.MultipartReader()
	if err != nil {
		s.Log.WithError(err).Error("error preparing multipart reader")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	created := make([]File, 0, 3)

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		contents, err := ioutil.ReadAll(part)
		if err != nil {
			s.Log.WithError(err).WithField(
				"filename", part.FileName(),
			).Error("error reading file part reader")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id := newID()
		// save contents to DB
		if err := writeBytes(s.DB, StorageFileKey, id, contents, meta, fg.TTL); err != nil {
			s.Log.WithError(err).Error("error writing file contents to store")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ct := part.Header.Get("Content-Type")
		if ct == "" || ct == "application/octet-stream" {
			if ext := filepath.Ext(part.FileName()); ext != "" {
				ct = mime.TypeByExtension(ext)
			} else {
				ct = http.DetectContentType(contents[:512])
			}
		}

		// save record to filegroup
		created = append(created, File{
			ID:       id,
			Name:     part.FormName(),
			FileName: part.FileName(),
			Mime:     ct,
		})
	}

	// re-fetch filegroup in case any changes happened
	if err := getOne(s.DB, DontBurn, StorageFileGroupKey, id, &fg); err != nil {
		if err == badger.ErrKeyNotFound {
			s.Log.WithField("id", id).Warn("FileGroup was deleted during file upload")
			for _, c := range created {
				if err := deleteRecord(s.DB, StorageFileKey, c.ID); err != nil {
					s.Log.WithError(err).WithField("fileID", c.ID).Error("while cleaning up file resources, got deletion error")
				}
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		s.Log.WithError(err).Error("error getting filegroup record")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fg.Files = append(fg.Files, created...)
	if err := writeType(s.DB, StorageFileGroupKey, fg.ID, fg, meta, fg.TTL); err != nil {
		s.Log.WithError(err).Error("error writing filegroup record")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

func (s *Server) FileGroupDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// delete files && group
	groupID := chi.URLParam(r, "id")

	var fg FileGroup
	if err := getOne(s.DB, DontBurn, StorageFileGroupKey, groupID, &fg); err != nil {
		if err == badger.ErrKeyNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.Log.WithError(err).Error("error getting filegroup record")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, f := range fg.Files {
		if err := deleteRecord(s.DB, StorageFileKey, f.ID); err != nil {
			if err == badger.ErrKeyNotFound {
				s.Log.WithFields(logrus.Fields{
					"groupID": groupID,
					"file":    f,
				}).Warn("file not found for deletion")
				continue
			}
			s.Log.WithField("fileID", f.ID).WithError(err).Error("error deleting file contents. May have dangling data")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if err := deleteRecord(s.DB, StorageFileGroupKey, groupID); err != nil {
		s.Log.WithField("groupID", groupID).WithError(err).Error("unable to delete file group")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) FileContentsGetHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "fid")

	contents, err := getOneBytes(s.DB, nil, StorageFileKey, id)
	if err == badger.ErrKeyNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		s.Log.WithError(err).WithField("id", id).Error("error fetching file contents")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")

	// if we have the group ID, we can try to get some metadata on the file
	if gid := chi.URLParam(r, "gid"); gid != "" {
		var fg FileGroup
		if err := getOne(s.DB, DontBurn, StorageFileGroupKey, gid, &fg); err == nil {
			for _, f := range fg.Files {
				if f.ID == id {
					if f.Mime != "" {
						w.Header().Set("Content-Type", f.Mime)
					}
					if r.URL.Query().Get("dl") != "" {
						w.Header().Set("Content-Disposition", `attachment; filename="`+f.FileName+`"`)
					}
					break
				}
			}
		}
	}

	w.Write(contents)
}

// DEBUG route for cleaning up of resources
func (s *Server) FileContentsListHandler(w http.ResponseWriter, r *http.Request) {
	infos, err := listInfo(s.DB, StorageFileKey)
	if err != nil {
		s.Log.WithError(err).Error("unable to get info on files")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsCfg.NewEncoder(w).Encode(map[string]interface{}{
		"data": infos,
	})
}

func (s *Server) FileCreateManualHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

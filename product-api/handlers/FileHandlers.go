package handlers

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/kjunn2000/go-server/files"
)

type Files struct {
	store files.Storage
}

func NewFileHandler(s files.Storage) *Files {
	return &Files{store: s}
}

func (f *Files) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]
	f.saveFile(id, fn, rw, r)
}

func (f *Files) saveFile(id, path string, rw http.ResponseWriter, r *http.Request) {
	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r.Body)
	if err != nil {
		log.Fatalln("Unable to save file", "error : ", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}

}

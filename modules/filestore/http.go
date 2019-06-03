package filestore

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"

	"github.com/spaceuptech/space-cloud/modules/auth"
)

// Supported content types
const (
	contentTypeApplicationJSON = "application/json"
	contentTypeEncodedForm     = "application/x-www-form-urlencoded"
	contentTypeMultiPartForm   = "multipart/form-data"
)

// HandleCreateFile creates the create file or directory endpoint
func (m *Module) HandleCreateFile(auth *auth.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Exit if file storage is not enabled
		if !m.isEnabled() {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "This feature isn't enabled"})
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Extract the path from the url
		token, project, _ := getMetaData(r)

		contentType := strings.Split(r.Header.Get("Content-type"), ";")[0]

		// Parse form
		var err error
		if contentType == contentTypeEncodedForm || contentType == contentTypeMultiPartForm {
			err = r.ParseMultipartForm((1 << 20) * 10)
		} else {
			err = r.ParseForm()
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
			return
		}

		path := r.Form.Get("path")
		fileType := r.Form.Get("fileType")
		makeAll, err := strconv.ParseBool(r.Form.Get("makeAll"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Incorrect value for makeAll"})
			return
		}

		// Check if the user is authorised to make this request
		err = auth.IsFileOpAuthorised(token, path, utils.FileCreate, map[string]interface{}{})
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to make this request"})
			return
		}

		if fileType == "file" {
			file, header, err := r.FormFile("file")
			defer file.Close()

			// Read file name from form if exists
			fileName := header.Filename
			if tempName := r.Form.Get("fileName"); tempName != "" {
				fileName = tempName
			}

			err = m.CreateFile(ctx, project, &model.CreateFileRequest{Name: fileName, Path: path, Type: fileType, MakeAll: makeAll}, file)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
		} else {
			name := r.Form.Get("name")
			err = m.CreateDir(ctx, project, &model.CreateFileRequest{Name: name, Path: path, Type: fileType, MakeAll: makeAll})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})
	}
}

// HandleRead creates read file and list directory endpoint
func (m *Module) HandleRead(auth *auth.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Exit if file storage is not enabled
		if !m.isEnabled() {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "This feature isn't enabled"})
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Extract the path from the url
		token, project, path := getMetaData(r)

		// Check if the user is authorised to make this request
		err := auth.IsFileOpAuthorised(token, path, utils.FileRead, map[string]interface{}{})
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to make this request"})
			return
		}

		op := r.URL.Query().Get("op")

		// List the specified directory if op type is list
		if op == "list" {
			mode := r.URL.Query().Get("mode")
			res, err := m.ListDir(ctx, project, &model.ListFilesRequest{Path: path, Type: mode})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": res})

			return
		}

		// Read the file from file storage
		file, err := m.ReadFile(ctx, project, path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		defer file.Close()

		w.WriteHeader(http.StatusOK)
		io.Copy(w, file.File)
	}
}

// HandleDelete creates read file and list directory endpoint
func (m *Module) HandleDelete(auth *auth.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Exit if file storage is not enabled
		if !m.isEnabled() {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "This feature isn't enabled"})
			return
		}

		// Create a context of execution
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Extract the path from the url
		token, project, path := getMetaData(r)

		// Check if the user is authorised to make this request
		err := auth.IsFileOpAuthorised(token, path, utils.FileDelete, map[string]interface{}{})
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to make this request"})
			return
		}

		err = m.DeleteDir(ctx, project, path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Give positive acknowledgement
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})
	}
}

func getMetaData(r *http.Request) (token string, project string, path string) {
	// Load the path parameters
	vars := mux.Vars(r)
	project = vars["project"]

	// Get the JWT token from header
	tokens, ok := r.Header["Authorization"]
	if !ok {
		tokens = []string{""}
	}
	token = strings.TrimPrefix(tokens[0], "Bearer ")
	a := strings.Split(r.URL.Path, "/")[5:]
	path = "/" + strings.Join(a, "/")
	return
}

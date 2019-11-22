package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/model"

	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/filestore"
)

// Supported content types
const (
	contentTypeApplicationJSON = "application/json"
	contentTypeEncodedForm     = "application/x-www-form-urlencoded"
	contentTypeMultiPartForm   = "multipart/form-data"
)

// HandleCreateFile creates the create file or directory endpoint
func HandleCreateFile(auth *auth.Module, fileStore *filestore.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the path from the url
		token, project, _ := getMetaData(r)
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

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

		if fileType == "file" {
			file, header, err := r.FormFile("file")
			defer file.Close()

			// Read file name from form if exists
			fileName := header.Filename
			if tempName := r.Form.Get("fileName"); tempName != "" {
				fileName = tempName
			}

			status, err := fileStore.UploadFile(ctx, project, token, &model.CreateFileRequest{Name: fileName, Path: path, Type: fileType, MakeAll: makeAll}, file)
			w.WriteHeader(status)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			json.NewEncoder(w).Encode(map[string]string{})
		} else {
			name := r.Form.Get("name")
			status, err := fileStore.CreateDir(ctx, project, token, &model.CreateFileRequest{Name: name, Path: path, Type: fileType, MakeAll: makeAll})
			w.WriteHeader(status)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			json.NewEncoder(w).Encode(map[string]string{})
		}
	}
}

// HandleRead creates read file and list directory endpoint
func HandleRead(auth *auth.Module, fileStore *filestore.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the path from the url
		token, project, path := getMetaData(r)
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		op := r.URL.Query().Get("op")

		// List the specified directory if op type is list
		if op == "list" {
			mode := r.URL.Query().Get("mode")

			status, res, err := fileStore.ListFiles(ctx, project, token, &model.ListFilesRequest{Path: path, Type: mode})
			w.WriteHeader(status)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{"result": res})
			return
		}

		// Read the file from file storage
		status, file, err := fileStore.DownloadFile(ctx, project, token, path)
		w.WriteHeader(status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		defer file.Close()
		io.Copy(w, file.File)
	}
}

// HandleDelete creates read file and list directory endpoint
func HandleDelete(auth *auth.Module, fileStore *filestore.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the path from the url
		token, project, path := getMetaData(r)
		defer r.Body.Close()

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		status, err := fileStore.DeleteFile(ctx, project, token, path)

		w.WriteHeader(status)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
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
	path = string(os.PathSeparator) + strings.Join(a, "/")
	return
}

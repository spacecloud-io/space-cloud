package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils/projects"
)

// Supported content types
const (
	contentTypeApplicationJSON = "application/json"
	contentTypeEncodedForm     = "application/x-www-form-urlencoded"
	contentTypeMultiPartForm   = "multipart/form-data"
)

// HandleCreateFile creates the create file or directory endpoint
func HandleCreateFile(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]

		defer r.Body.Close()

		// Load the project state
		state, err := projects.LoadProject(project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Exit if file storage is not enabled
		if !state.FileStore.IsEnabled() {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "This feature isn't enabled"})
			return
		}

		// Extract the path from the url
		token, project, _ := getMetaData(r)

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		contentType := strings.Split(r.Header.Get("Content-type"), ";")[0]
		if contentType == contentTypeEncodedForm || contentType == contentTypeMultiPartForm {
			err = r.ParseMultipartForm((1 << 20) * 10)
		} else {
			err = r.ParseForm()
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Could not parse form: %s", err.Error())})
			return
		}

		v := map[string]interface{}{}
		json.Unmarshal([]byte(r.FormValue("meta")), &v)
		path := r.FormValue("path")
		fileType := r.FormValue("fileType")
		var makeAll bool

		makeAllString := r.FormValue("makeAll")
		if makeAllString != "" {
			makeAll, err = strconv.ParseBool(makeAllString)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Incorrect value for makeAll"})
				return
			}
		}

		if fileType == "file" {
			file, header, err := r.FormFile("file")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Incorrect value for file: %s", err)})
				return
			}
			defer file.Close()

			// Read file name from form if exists
			fileName := header.Filename
			if tempName := r.FormValue("fileName"); tempName != "" {
				fileName = tempName
			}

			status, err := state.FileStore.UploadFile(ctx, project, token, &model.CreateFileRequest{Name: fileName, Path: path, Type: fileType, MakeAll: makeAll, Meta: v}, file)
			w.WriteHeader(status)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			json.NewEncoder(w).Encode(map[string]string{})
		} else {
			name := r.FormValue("name")
			status, err := state.FileStore.CreateDir(ctx, project, token, &model.CreateFileRequest{Name: name, Path: path, Type: fileType, MakeAll: makeAll}, v)
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
func HandleRead(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]

		// Extract the path from the url
		token, project, path := getMetaData(r)

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		op := r.URL.Query().Get("op")

		defer r.Body.Close()

		// Load the project state
		state, err := projects.LoadProject(project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// List the specified directory if op type is list
		if op == "list" {
			mode := r.URL.Query().Get("mode")

			status, res, err := state.FileStore.ListFiles(ctx, project, token, &model.ListFilesRequest{Path: path, Type: mode})
			w.WriteHeader(status)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{"result": res})
			return
		}

		// Read the file from file storage
		status, file, err := state.FileStore.DownloadFile(ctx, project, token, path)
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
func HandleDelete(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the path parameters
		vars := mux.Vars(r)
		project := vars["project"]

		// Extract the path from the url
		token, project, path := getMetaData(r)
		defer r.Body.Close()

		// Load the project state
		state, err := projects.LoadProject(project)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		v := map[string]interface{}{}
		json.NewDecoder(r.Body).Decode(v)

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		status, err := state.FileStore.DeleteFile(ctx, project, token, path, v)

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

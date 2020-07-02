package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Supported content types
const (
	// contentTypeApplicationJSON = "application/json"
	contentTypeEncodedForm   = "application/x-www-form-urlencoded"
	contentTypeMultiPartForm = "multipart/form-data"
)

// HandleCreateFile creates the create file or directory endpoint
func HandleCreateFile(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fileStore := modules.File()

		// Extract the path from the url
		token, projectID, _ := getFileStoreMeta(r)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		var err error
		contentType := strings.Split(r.Header.Get("Content-type"), ";")[0]
		if contentType == contentTypeEncodedForm || contentType == contentTypeMultiPartForm {
			err = r.ParseMultipartForm((1 << 20) * 10)
		} else {
			err = r.ParseForm()
		}
		if err != nil {
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Could not parse form: %s", err.Error()))
			return
		}

		v := map[string]interface{}{}
		_ = json.Unmarshal([]byte(r.FormValue("meta")), &v)
		path := r.FormValue("path")
		fileType := r.FormValue("fileType")
		var makeAll bool

		makeAllString := r.FormValue("makeAll")
		if makeAllString != "" {
			makeAll, err = strconv.ParseBool(makeAllString)
			if err != nil {
				_ = utils.SendErrorResponse(w, http.StatusBadRequest, "Incorrect value for makeAll")
				return
			}
		}

		if fileType == "file" {
			file, header, err := r.FormFile("file")
			if err != nil {
				_ = utils.SendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Incorrect value for file: %s", err))
				return
			}
			defer utils.CloseTheCloser(file)

			// Read file name from form if exists
			fileName := header.Filename
			if tempName := r.FormValue("fileName"); tempName != "" {
				fileName = tempName
			}

			status, err := fileStore.UploadFile(ctx, projectID, token, &model.CreateFileRequest{Name: fileName, Path: path, Type: fileType, MakeAll: makeAll, Meta: v}, file)
			if err != nil {
				_ = utils.SendErrorResponse(w, status, err.Error())
				return
			}
			_ = utils.SendResponse(w, status, map[string]string{})
		} else {
			name := r.FormValue("name")
			status, err := fileStore.CreateDir(ctx, projectID, token, &model.CreateFileRequest{Name: name, Path: path, Type: fileType, MakeAll: makeAll}, v)
			if err != nil {
				_ = utils.SendErrorResponse(w, status, err.Error())
				return
			}
			_ = utils.SendResponse(w, status, map[string]string{})
		}
	}
}

// HandleRead creates read file and list directory endpoint
func HandleRead(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fileStore := modules.File()

		// Extract the path from the url
		token, projectID, path := getFileStoreMeta(r)
		defer utils.CloseTheCloser(r.Body)

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		op := r.URL.Query().Get("op")

		// List the specified directory if op type is list
		if op == "list" {
			mode := r.URL.Query().Get("mode")
			status, res, err := fileStore.ListFiles(ctx, projectID, token, &model.ListFilesRequest{Path: path, Type: mode})
			if err != nil {
				_ = utils.SendErrorResponse(w, status, err.Error())
				return
			}
			_ = utils.SendResponse(w, status, map[string]interface{}{"result": res})
			return
		} else if op == "exist" {
			if err := fileStore.DoesExists(ctx, projectID, token, path); err != nil {
				_ = utils.SendErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}
			_ = utils.SendOkayResponse(w)
			return
		}

		// Read the file from file storage
		status, file, err := fileStore.DownloadFile(ctx, projectID, token, path)
		if err != nil {
			_ = utils.SendErrorResponse(w, status, err.Error())
			return
		}
		defer func() { _ = file.Close() }()
		w.WriteHeader(http.StatusOK)
		_, _ = io.Copy(w, file.File)
	}
}

// HandleDelete creates read file and list directory endpoint
func HandleDelete(modules *modules.Modules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fileStore := modules.File()

		// Extract the path from the url
		token, projectID, path := getFileStoreMeta(r)
		defer utils.CloseTheCloser(r.Body)

		v := map[string]interface{}{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		log.Println("v", v)

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		status, err := fileStore.DeleteFile(ctx, projectID, token, path, v)

		if err != nil {
			_ = utils.SendErrorResponse(w, status, err.Error())
			return
		}
		_ = utils.SendResponse(w, status, map[string]string{})
	}
}

func getFileStoreMeta(r *http.Request) (token string, projectID string, path string) {
	// Load the path parameters
	vars := mux.Vars(r)
	projectID = vars["project"]

	// Get the JWT token from header
	tokens, ok := r.Header["Authorization"]
	if !ok {
		tokens = []string{""}
	}
	token = strings.TrimPrefix(tokens[0], "Bearer ")
	a := strings.Split(r.URL.Path, "/")
	for index, value := range a {
		if value == "files" {
			path = string(os.PathSeparator) + strings.Join(a[index+1:], "/")
		}
	}
	return
}

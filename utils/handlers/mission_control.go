package handlers

import (
	"net/http"
	"os"
	"strings"
)

// HandleMissionControl hosts the static resources for mission control
func HandleMissionControl(staticPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path

		defer r.Body.Close()

		path := strings.TrimPrefix(url, "/mission-control")
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		path = staticPath + path

		// Check if path exists
		if fileInfo, err := os.Stat(path); !os.IsNotExist(err) {
			// If path exists and is of type file then serve that file
			if !fileInfo.IsDir() {
				http.ServeFile(w, r, path)
				return
			}
			// Else if a index file exists within that folder serve that index file
			path = strings.TrimSuffix(path, "/")
			if _, err := os.Stat(path + "/index.html"); !os.IsNotExist(err) {
				http.ServeFile(w, r, path+"/index.html")
				return
			}
		}

		// If path does not exists serve the root index file
		http.ServeFile(w, r, strings.TrimSuffix(staticPath, "/")+"/index.html")
	}
}

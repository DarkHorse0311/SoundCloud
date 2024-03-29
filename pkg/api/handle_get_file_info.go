package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
)

type GetFileRequest struct {
	ID int64 `json:"id"`
}

func (api *API) HandleGetFileInfo(w http.ResponseWriter, r *http.Request) {
	getFileRequest := &GetFileRequest{
		ID: -1,
	}

	err := json.NewDecoder(r.Body).Decode(getFileRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check empty
	if getFileRequest.ID < 0 {
		api.HandleErrorString(w, r, `"id" can't be none or negative`)
		return
	}

	log.Println("[api] Get file info", getFileRequest.ID)

	file, err := api.Db.GetFile(getFileRequest.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	err = json.NewEncoder(w).Encode(file)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}

func (api *API) HandleGetFileFfprobeInfo(w http.ResponseWriter, r *http.Request) {
	getFileRequest := &GetFileRequest {
		ID: -1,
	}

	err := json.NewDecoder(r.Body).Decode(getFileRequest)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// check empty
	if getFileRequest.ID < 0 {
		api.HandleErrorString(w, r, `"id" can't be none or negative`)
		return
	}

	log.Println("[api] Get file Ffprobe info", getFileRequest.ID)

	file, err := api.Db.GetFile(getFileRequest.ID)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	path, err := file.Path()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	cmd := exec.Command("ffprobe", "-i", path, "-hide_banner")
	cmd.Stderr = w
	err = cmd.Run()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
}

// /api/v1/get_file?id=123
// get raw file with io.Copy method
func (api *API) HandleGetFile(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ids := q["id"]
	_id, err := strconv.Atoi(ids[0])
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	id := int64(_id)

	// check empty
	if id < 0 {
		api.HandleErrorString(w, r, `"id" can't be none or negative`)
		return
	}

	file, err := api.Db.GetFile(id)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	path, err := file.Path()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	log.Println("[api] Get pipe raw file", path)

	src, err := os.Open(path)
	if err != nil {
		api.HandleError(w, r, err)
		return
	}
	defer src.Close()
	io.Copy(w, src)
}

// /get_file_direct?id=1
// get raw file with http.ServeFile method
func (api *API) HandleGetFileDirect(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ids := q["id"]
	if len(ids) == 0 {
		api.HandleErrorString(w, r, `parameter "id" can't be empty`)
		return
	}
	id, err := strconv.Atoi(ids[0])
	if err != nil {
		api.HandleErrorString(w, r, `parameter "id" should be an integer`)
		return
	}
	file, err := api.Db.GetFile(int64(id))
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	path, err := file.Path()
	if err != nil {
		api.HandleError(w, r, err)
		return
	}

	// set header for filename
	filename := file.Filename
	// encode filename to URL
	filename = url.PathEscape(filename)
	w.Header().Set("Content-Disposition", "inline; filename*=UTF-8''"+filename)

	log.Println("[api] Get direct raw file", path)

	http.ServeFile(w, r, path)
}

package chat

import (
	"encoding/json"
	"net/http"
	"strconv"

	errors "github.com/Trecer05/Swiftly/internal/errors/global"
	fileManager "github.com/Trecer05/Swiftly/internal/filemanager"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"

	"github.com/gorilla/mux"
)

func UploadImgHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, t models.ChatType) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	photos := r.MultipartForm.File["photos"]
	if len(photos) == 0 {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrNoPhotos, http.StatusBadRequest)
	}

	urls, err := fileManager.SaveMessageFiles(photos, id, t, models.DataTypeImg)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}

func UploadVideoHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, t models.ChatType) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	videos := r.MultipartForm.File["videos"]
	if len(videos) == 0 {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrNoPhotos, http.StatusBadRequest)
	}

	urls, err := fileManager.SaveMessageFiles(videos, id, t, models.DataTypeVid)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}

func UploadAudioHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, t models.ChatType) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	audio := r.MultipartForm.File["audio"]
	if len(audio) == 0 {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrNoPhotos, http.StatusBadRequest)
	}

	urls, err := fileManager.SaveMessageFiles(audio, id, t, models.DataTypeAud)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, t models.ChatType) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrNoPhotos, http.StatusBadRequest)
	}

	urls, err := fileManager.SaveMessageFiles(files, id, t, models.DataTypeDoc)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}

func UploadImgVideoHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, t models.ChatType) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrNoPhotos, http.StatusBadRequest)
	}

	urls, err := fileManager.SaveMessageFiles(files, id, t, models.DataTypeImgVid)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}

func GetFilesHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, t models.ChatType) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	filesUrls, err := fileManager.GetMedias(id, t)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filesUrls)
}

func GetImgHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, t models.ChatType) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	url := vars["url"]

	file, fileType, err := fileManager.GetFile(url, id, t, models.DataTypeImg)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", fileType)
	if _, err := w.Write(file); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
}

func GetVideoHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, t models.ChatType) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	url := vars["url"]

	file, fileType, err := fileManager.GetFile(url, id, t, models.DataTypeVid)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", fileType)
	if _, err := w.Write(file); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
}

func GetAudioHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, t models.ChatType) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	url := vars["url"]

	file, fileType, err := fileManager.GetFile(url, id, t, models.DataTypeAud)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", fileType)
	if _, err := w.Write(file); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
}

func GetFileHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, t models.ChatType) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	url := vars["url"]

	file, fileType, err := fileManager.GetFile(url, id, t, models.DataTypeDoc)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", fileType)
	if _, err := w.Write(file); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
}

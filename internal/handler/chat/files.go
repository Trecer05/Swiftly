package chat

import (
	"encoding/json"
	"net/http"
	"strconv"

	errors "github.com/Trecer05/Swiftly/internal/errors/global"
	authErrors "github.com/Trecer05/Swiftly/internal/errors/auth"
	fileManager "github.com/Trecer05/Swiftly/internal/filemanager"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
	service "github.com/Trecer05/Swiftly/internal/service/chat"
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

func GetProfileAvatarUrlsHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewErrorBody(w, "application/json", authErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	names, err := fileManager.GetUserPhotosUrls(userId)
	if err != nil {
		if names == nil {
			serviceHttp.NewErrorBody(w, "application/json", errors.ErrNoPhotos, http.StatusNotFound)
			return
		}
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{
		"avatars": names,
	})
}

func GetProfileAvatarHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewErrorBody(w, "application/json", authErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	url := vars["url"]

	file, fileType, err := fileManager.GetUserPhotoByUrl(userId, url)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	if file == nil {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrNoPhotos, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", fileType)
	if _, err := w.Write(file); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
}

func GetUserAvatarUrlsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := service.GetIdFromVars(r)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	names, err := fileManager.GetUserPhotosUrls(id)
	if err != nil {
		if names == nil {
			serviceHttp.NewErrorBody(w, "application/json", errors.ErrNoPhotos, http.StatusNotFound)
			return
		}
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{
		"avatars": names,
	})
}

func GetUserAvatarHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	url := vars["url"]

	file, fileType, err := fileManager.GetUserPhotoByUrl(id, url)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	if file == nil {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrNoPhotos, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", fileType)
	if _, err := w.Write(file); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
}

func UploadProfileAvatarHandler(w http.ResponseWriter, r *http.Request) {
	if err := fileManager.AddUserPhoto(r, r.Context().Value("id").(int)); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func DeleteProfileAvatarHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userId := r.Context().Value("id").(int)
	if err := fileManager.DeleteUserAvatar(vars["url"], userId); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func GetGroupAvatarUrlHandler(w http.ResponseWriter, r *http.Request) {
	id, err := service.GetIdFromVars(r)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	url, err := fileManager.GetLatestGroupAvatarUrl(id)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	if url == "" {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrNoPhotos, http.StatusNotFound)
		return
	} 

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"avatar": url,
	})
}

func GetGroupAvatarHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	url := vars["url"]

	file, fileType, err := fileManager.GetGroupPhotoByUrl(id, url)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	if file == nil {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrNoPhotos, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", fileType)
	if _, err := w.Write(file); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
}

func UploadGroupAvatarHandler(w http.ResponseWriter, r *http.Request) {
	id, err := service.GetIdFromVars(r)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	if err := fileManager.AddGroupPhoto(r, id); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func DeleteGroupAvatarHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	if err := fileManager.DeleteGroupAvatar(vars["url"], id); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

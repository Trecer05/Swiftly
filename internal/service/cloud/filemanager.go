package cloud

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	cloudFilemanager "github.com/Trecer05/Swiftly/internal/filemanager/cloud"
	chatModels "github.com/Trecer05/Swiftly/internal/model/chat"
	models "github.com/Trecer05/Swiftly/internal/model/cloud"
	kafkaModels "github.com/Trecer05/Swiftly/internal/model/kafka"
	cloudKafkaManager "github.com/Trecer05/Swiftly/internal/repository/kafka/cloud"
	postgres "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"

	errors "github.com/Trecer05/Swiftly/internal/errors/file"
)

// GetFile checks permissions and returns file model and its data.
// requestUserId is the ID of the user making the request
func GetFileSync(fileModel *models.File, requestUserID int, kafkaManager *cloudKafkaManager.KafkaManager) ([]byte, error) {
	if err := checkAccess(fileModel, requestUserID, kafkaManager); err != nil {
		return nil, err
	}

	data, _, err := cloudFilemanager.GetFileSync(fileModel)
	return data, err
}

func GetFileStream(fileModel *models.File, requestUserID int, kafkaManager *cloudKafkaManager.KafkaManager) (*os.File, error) {
	if err := checkAccess(fileModel, requestUserID, kafkaManager); err != nil {
		return nil, err
	}
	return cloudFilemanager.GetFileStream(fileModel)
}

func checkAccess(file *models.File, requestUserID int, kafkaManager *cloudKafkaManager.KafkaManager) error {
	switch file.OwnerType {
	case models.OwnerTypeUser:
		switch file.Visibility {
		case models.VisibilityPrivate:
			if file.OwnerID != requestUserID {
				return errors.ErrPermissionDenied
			}
		case models.VisibilityShared:
			if err := CheckUserInTeam(file.OwnerID, requestUserID, kafkaManager); err != nil {
				return err
			}

		}
	case models.OwnerTypeTeam:
		switch file.Visibility {
		case models.VisibilityPublic:
			if err := CheckUserInTeam(file.OwnerID, requestUserID, kafkaManager); err != nil {
				return err
			}
		case models.VisibilityPrivate:
			if file.OwnerID != requestUserID {
				return errors.ErrPermissionDenied
			}
		case models.VisibilityShared:
			return nil
		}
	}
	return errors.ErrPermissionDenied // по умолчанию лучше верну запрет, лучше лишний раз отказать чем разрешить
}

// Новая функция для более оптимизированной проверки доступа при получении папки команды
// чтобы каждый раз не обращаться к кафке для каждой проверки файла
// принимаем bool значение, состоит ли пользователь в организации
func HasAccessToTeamFile(file *models.File, requestUserID int, isInTeam bool) error {
	switch file.OwnerType {
	case models.OwnerTypeUser:
		switch file.Visibility {
		case models.VisibilityPrivate:
			if file.OwnerID != requestUserID {
				return errors.ErrPermissionDenied
			}
		case models.VisibilityShared:
			if !isInTeam {
				return errors.ErrPermissionDenied
			}

		}
	case models.OwnerTypeTeam:
		switch file.Visibility {
		case models.VisibilityPublic:
			if !isInTeam {
				return errors.ErrPermissionDenied
			}
		case models.VisibilityPrivate:
			if file.OwnerID != requestUserID {
				return errors.ErrPermissionDenied
			}
		case models.VisibilityShared:
			return nil
		}
	}
	return errors.ErrPermissionDenied
}

func CheckUserInTeam(teamID int, requestUserID int, kafkaManager *cloudKafkaManager.KafkaManager) error {
	corrID, err := kafkaManager.SendMessage(context.Background(), "check_user", kafkaModels.CheckUserInTeam{TeamID: teamID, UserID: requestUserID})
	if err != nil {
		return err
	}

	resp, err := kafkaManager.WaitForResponse(corrID.String(), 5*time.Second)
	if err != nil {
		return err
	}

	var r chatModels.CheckUserInTeamResponse
	if err := json.Unmarshal(resp.Payload, &r); err != nil {
		logger.Logger.Errorf("Error unmarshaling check user in team response: %v", err)
		return err
	}
	if !r.IsInTeam {
		return errors.ErrPermissionDenied
	}
	return nil
}

func UpdateUserFileNameByID(mgr *postgres.Manager, fileID, newFilename string, userID int) (time.Time, error) {
	filepath, err := mgr.GetUserFilepathByID(userID, fileID)
	if err != nil {
		return time.Time{}, err
	}

	newOrigFilename, newFilepath, err := cloudFilemanager.UpdateFileName(filepath, newFilename)
	if err != nil {
		return time.Time{}, err
	}

	return mgr.UpdateFileFilenameByID(userID, fileID, newOrigFilename, newFilename, newFilepath)
}

func UpdateTeamFileNameByID(mgr *postgres.Manager, fileID, newFilename string, teamID int, userID int) (time.Time, error) {
	filepath, err := mgr.GetTeamFilepathByID(teamID, fileID)
	if err != nil {
		return time.Time{}, err
	}

	newOrigFilename, newFilepath, err := cloudFilemanager.UpdateFileName(filepath, newFilename)
	if err != nil {
		return time.Time{}, err
	}

	return mgr.UpdateFileFilenameByID(userID, fileID, newOrigFilename, newFilename, newFilepath)
}

func UpdateUserFolderNameByID(mgr *postgres.Manager, folderID, newFoldername string, userID int) (time.Time, error) {
	filepath, err := mgr.GetUserFolderpathByID(userID, folderID)
	if err != nil {
		return time.Time{}, err
	}

	newFolderName, newFolderpath, err := cloudFilemanager.UpdateFolderName(filepath, newFoldername)
	if err != nil {
		return time.Time{}, err
	}

	return mgr.UpdateFolderFoldernameByID(userID, folderID, newFolderName, newFolderpath)
}

func UpdateTeamFolderNameByID(mgr *postgres.Manager, folderID, newFoldername string, teamID int, userID int) (time.Time, error) {
	filepath, err := mgr.GetTeamFolderpathByID(userID, folderID)
	if err != nil {
		return time.Time{}, err
	}

	newFolderName, newFolderpath, err := cloudFilemanager.UpdateFolderName(filepath, newFoldername)
	if err != nil {
		return time.Time{}, err
	}

	return mgr.UpdateFolderFoldernameByID(userID, folderID, newFolderName, newFolderpath)
}

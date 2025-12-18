package cloud

import (
	"os"

	cloudFilemanager "github.com/Trecer05/Swiftly/internal/filemanager/cloud"
	models "github.com/Trecer05/Swiftly/internal/model/cloud"
	chatManager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"

	errors "github.com/Trecer05/Swiftly/internal/errors/file"
)

// GetFile checks permissions and returns file model and its data.
// requestUserId is the ID of the user making the request
func GetFileSync(fileModel *models.File, requestUserID int, chatManager *chatManager.Manager) ([]byte, error) {
	if err := checkAccess(fileModel, requestUserID, chatManager); err != nil {
		return nil, err
	}

	data, _, err := cloudFilemanager.GetFileSync(fileModel)
	return data, err
}

func GetFileStream(fileModel *models.File, requestUserID int, chatManager *chatManager.Manager) (*os.File, error) {
	if err := checkAccess(fileModel, requestUserID, chatManager); err != nil {
		return nil, err
	}
	return cloudFilemanager.GetFileStream(fileModel)
}

func checkAccess(file *models.File, requestUserID int, chatManager *chatManager.Manager) error {
	switch file.Visibility {
	case models.VisibilityPublic:
		switch file.OwnerType {
		case models.OwnerTypeUser:
			if file.OwnerID != requestUserID {
				return errors.ErrPermissionDenied
			}
		case models.OwnerTypeTeam:
			if isExist, err := chatManager.IsUserInTeam(file.OwnerID, requestUserID); err != nil || !isExist {
				return err
			}
		default:
			return errors.ErrPermissionDenied
		}
	case models.VisibilityPrivate:
		if file.OwnerID != requestUserID {
			return errors.ErrPermissionDenied
		}

	case models.VisibilityShared:
		// доступ разрешён
	}
	return nil
}

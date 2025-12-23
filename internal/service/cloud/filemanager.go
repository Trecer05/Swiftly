package cloud

import (
	"context"
	"encoding/json"
	"os"
	"time"

	cloudFilemanager "github.com/Trecer05/Swiftly/internal/filemanager/cloud"
	chatModels "github.com/Trecer05/Swiftly/internal/model/chat"
	models "github.com/Trecer05/Swiftly/internal/model/cloud"
	kafkaModels "github.com/Trecer05/Swiftly/internal/model/kafka"
	cloudKafkaManager "github.com/Trecer05/Swiftly/internal/repository/kafka/cloud"

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
			if err := checkUserInTeam(file, requestUserID, kafkaManager); err != nil {
				return err
			}

		}
	case models.OwnerTypeTeam:
		switch file.Visibility {
		case models.VisibilityPublic:
			if err := checkUserInTeam(file, requestUserID, kafkaManager); err != nil {
				return err
			}
		case models.VisibilityPrivate:
			if file.OwnerID != requestUserID {
				return errors.ErrPermissionDenied
			}
		case models.VisibilityShared:
			return nil
		}
		// case models.OwnerTypeTeam:
		// 	switch file.Visibility {
		// 	case models.VisibilityPublic:
		// 		isExist, err := kafkaManager.IsUserInTeam(file.OwnerID, requestUserID)
		// 		if err != nil || !isExist {
		// 		}

		// }
	}
	return errors.ErrPermissionDenied // по умолчанию лучше верну запрет, лучше лишний раз отказать чем разрешить
}

func checkUserInTeam(file *models.File, requestUserID int, kafkaManager *cloudKafkaManager.KafkaManager) error {
	corrID, err := kafkaManager.SendMessage(context.Background(), "check_user", kafkaModels.CheckUserInTeam{TeamID: file.OwnerID, UserID: requestUserID})
	if err != nil {
		return err
	}

	resp, err := kafkaManager.WaitForResponse(corrID.String(), 5*time.Second)
	if err != nil {
		return err
	}

	var r chatModels.CheckUserInTeamResponse
	if err := json.Unmarshal(resp.Payload, &r); err != nil {
		return err
	}
	if !r.IsInTeam {
		return errors.ErrPermissionDenied
	}
	return nil
}

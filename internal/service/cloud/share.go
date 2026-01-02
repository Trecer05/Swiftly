package cloud

import (
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"
	"github.com/google/uuid"
)

func ShareFile(fileUUID uuid.UUID, teamID int, userID int, mgr *manager.Manager) (string, error) {
	if err := mgr.ShareTeamFileByID(fileUUID.String(), teamID, userID); err != nil {
		return "", err
	}

	return GenerateShareFileLink(fileUUID.String()), nil
}

func ShareFolder(folderUUID uuid.UUID, teamID int, userID int, mgr *manager.Manager) (string, error) {
}

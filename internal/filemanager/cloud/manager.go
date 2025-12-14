package cloud

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/Trecer05/Swiftly/internal/config/logger"
)

func CreateStartDirs() {
	commandsDir := filepath.Join("cloud", "commands")
	usersDir := filepath.Join("cloud", "users")

	if err := createDirIfNotExist(commandsDir); err != nil {
		logger.Logger.Fatal("Failed to create commands directory: ", err)
	}

	if err := createDirIfNotExist(usersDir); err != nil {
		logger.Logger.Fatal("Failed to create users directory: ", err)
	}
}

func createDirIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	return nil
}

func CreateTeamFolders(teamID int) error {
	if err := createDirIfNotExist(filepath.Join("cloud", "teams", strconv.Itoa(teamID))); err != nil {
		return err
	}

	return nil
}

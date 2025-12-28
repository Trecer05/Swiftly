package cloud

import "os"

func DeleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}

func DeleteFolder(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}

	return nil
}
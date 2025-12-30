package cloud

import "fmt"

func GenerateShareFileLink(fileID string) string {
	return fmt.Sprintf("/sharedfiles/file/%s", fileID)
}
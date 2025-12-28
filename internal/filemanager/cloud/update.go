package cloud

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func UpdateFileName(filePath, newFilename string) (string, string, error) {
    dir := filepath.Dir(filePath)
    originalExt := filepath.Ext(filePath)
    
    newName := newFilename
    
    timestamp := fmt.Sprintf("%d_", time.Now().UnixNano())
    newName = timestamp + newName
    
    if originalExt != "" && filepath.Ext(newName) == "" {
        newName += originalExt
    }
    
    newPath := filepath.Join(dir, newName)
    
    if _, err := os.Stat(newPath); err == nil {
        newName = fmt.Sprintf("%s_%d", newName[:len(newName)-len(filepath.Ext(newName))], time.Now().UnixNano())
        if originalExt != "" {
            newName += originalExt
        }
        newPath = filepath.Join(dir, newName)
    }
    
    if err := os.Rename(filePath, newPath); err != nil {
        return "", "", err
    }
    
    return newName, newPath, nil
}

func UpdateFolderName(folderPath, newFolderName string) (string, string, error) {
    fileInfo, err := os.Stat(folderPath)
    if err != nil {
        if os.IsNotExist(err) {
            return "", "", err
        }
        return "", "", err
    }
    
    if !fileInfo.IsDir() {
        return "", "", err
    }
    
    parentDir := filepath.Dir(folderPath)
    
    cleanedName := cleanFolderName(newFolderName)
    if cleanedName == "" {
        return "", "", err
    }
    
    newPath := filepath.Join(parentDir, cleanedName)
    
    if _, err := os.Stat(newPath); err == nil {
        timestamp := fmt.Sprintf("_%d", time.Now().UnixNano())
        cleanedName += timestamp
        newPath = filepath.Join(parentDir, cleanedName)
    }
    
    if err := os.Rename(folderPath, newPath); err != nil {
        return "", "", err
    }
    
    return cleanedName, newPath, nil
}

func cleanFolderName(name string) string {
    name = strings.TrimSpace(name)
    
    invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*", "\\", "/", "\x00"}
    
    for _, char := range invalidChars {
        name = strings.ReplaceAll(name, char, "_")
    }
    
    name = strings.Trim(name, ".")
    
    reservedNames := []string{"CON", "PRN", "AUX", "NUL", 
        "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
        "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
    
    upperName := strings.ToUpper(name)
    for _, reserved := range reservedNames {
        if upperName == reserved || strings.HasPrefix(upperName, reserved+".") {
            return "_" + name
        }
    }
    
    return name
}

package files

import (
	"os"
)

const CONFIG_PATH = "/.config/moon/"
const AUTH_FILENAME = "auth"

func saveToFile(path string, content []byte) error {
	return os.WriteFile(path, content, 0644)
}
func readFromFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func SaveToConfigFile(filename string, content []byte) error {
	homePath := os.Getenv("HOME")
	return saveToFile(homePath+CONFIG_PATH+filename, content)
}
func ReadFromConfigFile(filename string) ([]byte, error) {
	homePath := os.Getenv("HOME")
	return readFromFile(homePath + CONFIG_PATH + filename)
}
func IsFromConfigFileExist(filename string) bool {
	homePath := os.Getenv("HOME")
	if _, err := os.Stat(homePath + CONFIG_PATH + filename); err == nil {
		return true
	} 
	return false
}

func InitConfigFolders() error {
	homePath := os.Getenv("HOME")
	return os.MkdirAll(homePath+CONFIG_PATH, 0755)
}

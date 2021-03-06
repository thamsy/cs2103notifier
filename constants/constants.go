package constants

import "os"

const (
	// Messages
	START_MESSAGE = "CS2103 Notifier Started..."
	START_LOG     = "Log Starting..."
	END_MESSAGE   = "CS2103 Notifier Closing..."

	// Directory Names
	STORAGE_FOLDER = ".cs2103notifier"
	PREV_PREFIX    = "prev_"
	LOG_FILENAME   = "cs2103notifierlog.txt"
)

// Get Directory
func GetStorageDir() string {
	return os.Getenv("HOME") + "/" + STORAGE_FOLDER
}

func GetCurrentDir(filename string) string {
	return GetStorageDir() + "/" + filename
}

func GetPrevDir(filename string) string {
	return GetStorageDir() + "/" + PREV_PREFIX + filename
}

func GetLogDir() string {
	return GetStorageDir() + "/" + LOG_FILENAME
}

package logging

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
)

// GenerateID returns a random 16-byte hex string (32 hex chars).
// We use crypto/rand for strong uniqueness properties.
func GenerateID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		// Unlikely path — return zeros if entropy source fails.
		return hex.EncodeToString(b[:])
	}
	return hex.EncodeToString(b[:])
}

// CreateAppDataFolder creates an application data folder in the user's cache directory.
func CreateAppDataFolder(applicationName string) (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir = dir + "\\" + applicationName
	err = os.MkdirAll(dir, 0600)
	if err != nil {
		return "", err
	}
	return dir, nil
}

// OpenLogFile opens (or creates) a log file for appending log entries.
func OpenLogFile(fileName string) (*os.File, error) {
	// open the log file for appending log entries
	// create it if it does not exist with permissions rw-r--r--
	// append mode so we do not overwrite existing logs
	fi, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// log file not ready so default std.err logging here
		slog.Error(fmt.Sprintf("%s\n", "Failed to create logfile for writing"))
		slog.Error(err.Error())
		return &os.File{}, err
	}
	return fi, nil
}

func LoggerOptions() slog.HandlerOptions {
	// TODO: adjust options based on environment
	var options slog.HandlerOptions
	options = slog.HandlerOptions{AddSource: false}
	return options
}
